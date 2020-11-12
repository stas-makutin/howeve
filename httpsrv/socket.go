package httpsrv

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/stas-makutin/howeve/events"
	"github.com/stas-makutin/howeve/events/handlers"
	"github.com/stas-makutin/howeve/log"
)

const wsopStart = "S"
const wsopFinish = "F"
const wsopOutbound = "O"
const wsopInbound = "I"

const wsocSuccess = "0"
const wsocReadError = "1"
const wsocWriteError = "2"
const wsocUnexpectedRequest = "3"
const wsocUnexpectedResponse = "4"
const wsocNoRequestMapping = "5"

var connectionOrdinal handlers.Ordinal

func handleWebsocket(w http.ResponseWriter, r *http.Request, hc *handlerContext) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		appendLogFields(r, fmt.Sprintf("%T: %v", err, err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	co := connectionOrdinal.Next()
	log.Report(log.SrcWS, wsopStart, co.String())

	closeCh := make(chan struct{})
	go func() {
		select {
		case <-hc.stopCh: // closes connection on http server shutdown
		case <-closeCh: // closes connection if it is terminated
		}
		conn.Close()
		log.Report(log.SrcWS, wsopFinish, co.String())
	}()

	hc.handlerWg.Add(1)
	go func() {
		defer hc.handlerWg.Done()
		defer close(closeCh)
		messageLoop(conn, co)
	}()
}

func messageLoop(conn net.Conn, co handlers.Ordinal) {
	var id events.SubscriberID
	id = handlers.Dispatcher.Subscribe(func(event interface{}) {
		if te, ok := event.(events.TargetedResponse); ok && te.Receiver() == id {
			var eo handlers.Ordinal
			var eid string
			if th, ok := event.(handlers.TraceHeader); ok {
				eo, eid = th.Ordinal(), th.TraceID()
			}
			if query := queryFromEvent(te); query != nil {
				writeQuery(conn, co, eo, query)
			} else {
				log.Report(log.SrcWS, wsopOutbound, co.String(), eo.String(), wsocUnexpectedResponse, eid, fmt.Sprintf("%T", event))
			}
		}
	})
	defer handlers.Dispatcher.Unsubscribe(id)

	for {
		msg, _, err := wsutil.ReadClientData(conn)
		if err != nil {
			if err != io.EOF {
				log.Report(log.SrcWS, wsopInbound, co.String(), handlers.EventOrdinal.Next().String(), wsocReadError, err.Error())
			}
			break
		}

		var q Query
		err = json.Unmarshal(msg, &q)
		if err != nil {
			sg := msg
			if len(sg) > 32 {
				sg = sg[:32]
			}
			eo := handlers.EventOrdinal.Next()
			log.Report(log.SrcWS, wsopInbound, co.String(), eo.String(), wsocUnexpectedRequest, string(sg), err.Error())
			writeQuery(conn, co, eo, &Query{Type: queryUnexpected, ID: q.ID})
			continue
		}

		// process message
		n, _ := queryNameMap[q.Type]
		if event := q.toTargetedRequest(id); event != nil {
			var eo handlers.Ordinal
			if th, ok := event.(handlers.TraceHeader); ok {
				eo = th.Ordinal()
			}
			log.Report(log.SrcWS, wsopInbound, co.String(), eo.String(), wsocSuccess, q.ID, n, strconv.FormatInt(int64(len(msg)), 10))
			handlers.Dispatcher.Send(event)
		} else {
			eo := handlers.EventOrdinal.Next()
			log.Report(log.SrcWS, wsopInbound, co.String(), eo.String(), wsocNoRequestMapping, q.ID, n, strconv.FormatInt(int64(len(msg)), 10))
			writeQuery(conn, co, eo, &Query{Type: queryUnexpected, ID: q.ID})
		}
	}
}

func writeQuery(conn net.Conn, co handlers.Ordinal, eo handlers.Ordinal, q *Query) {
	w := newWebSocketTextWriter(conn)
	n, _ := queryNameMap[q.Type]
	if err := json.NewEncoder(w).Encode(q); err == nil {
		log.Report(log.SrcWS, wsopOutbound, co.String(), eo.String(), wsocSuccess, q.ID, n, strconv.FormatUint(w.length, 10))
	} else {
		log.Report(log.SrcWS, wsopOutbound, co.String(), eo.String(), wsocWriteError, q.ID, n, err.Error())
	}
}

// WebSocketTextWriter struct
type WebSocketTextWriter struct {
	conn   net.Conn
	length uint64
}

func newWebSocketTextWriter(conn net.Conn) *WebSocketTextWriter {
	return &WebSocketTextWriter{conn, 0}
}

func (w *WebSocketTextWriter) Write(p []byte) (n int, err error) {
	err = wsutil.WriteServerText(w.conn, p)
	if err == nil {
		n = len(p)
		w.length += uint64(n)
	}
	return
}
