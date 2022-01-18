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
	"golang.org/x/net/context"
)

const wsopStart = "S"
const wsopFinish = "F"
const wsopOutbound = "O"
const wsopInbound = "I"

const wsocSuccess = "0"
const wsocReadError = "R"
const wsocWriteError = "W"
const wsocUnexpectedRequest = "R"
const wsocUnexpectedResponse = "P"
const wsocNoRequestMapping = "M"
const wsocSkipResponse = "S"

var connectionOrdinal handlers.Ordinal

func handleWebsocket(w http.ResponseWriter, r *http.Request, hc *handlerContext) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		appendLogFields(r, fmt.Sprintf("%T: %v", err, err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	hc.handlerWg.Add(1)
	go func() {
		defer hc.handlerWg.Done()
		messageLoop(conn, hc.stopCh)
	}()
}

func messageLoop(conn net.Conn, stopCh chan struct{}) {
	co := connectionOrdinal.Next()
	log.Report(log.SrcWS, wsopStart, co.String())
	defer log.Report(log.SrcWS, wsopFinish, co.String())
	defer conn.Close()

	writeCh := make(chan interface{}, 32)
	subscriptions := newSocketSubscription()

	var id events.SubscriberID
	id = handlers.Dispatcher.Subscribe(func(event interface{}) {
		toWrite := false
		if subscriptions.subscribed(event) {
			toWrite = true
		} else if te, ok := event.(events.TargetedResponse); ok && te.Receiver() == id {
			toWrite = true
		}
		if toWrite {
		Written:
			for {
				select {
				case writeCh <- event:
					break Written
				default:
					skippedEvent := <-writeCh
					var eo handlers.Ordinal
					var eid string
					if ti, ok := skippedEvent.(handlers.TraceInfo); ok {
						eo, eid = ti.Ordinal(), ti.TraceID()
					}
					log.Report(log.SrcWS, wsopOutbound, co.String(), eo.String(), wsocSkipResponse, eid, fmt.Sprintf("%T", skippedEvent))
				}
			}
		}
	})
	defer handlers.Dispatcher.Unsubscribe(id)

	closeCh := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())

	// read loop
	go func() {
		for {
			msg, _, err := wsutil.ReadClientData(conn)
			if err != nil {
				if !isConnClosed(err) {
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
				if writeQuery(conn, co, eo, &Query{Type: queryUnexpected, ID: q.ID}) {
					break
				}
				continue
			}

			// process message
			n := queryNameMap[q.Type]
			if q.Type == queryEventSubscribe {
				eo := handlers.EventOrdinal.Next()
				log.Report(log.SrcWS, wsopInbound, co.String(), eo.String(), wsocSuccess, q.ID, n, strconv.FormatInt(int64(len(msg)), 10))
				subscriptions.subscribe(q.Payload.(*Subscription))
				eo = handlers.EventOrdinal.Next()
				if writeQuery(conn, co, eo, &Query{Type: queryEventSubscribeResult, ID: q.ID}) {
					break
				}
			} else if event := q.toTargetedRequest(ctx, id); event != nil {
				var eo handlers.Ordinal
				if ti, ok := event.(handlers.TraceInfo); ok {
					eo = ti.Ordinal()
				}
				log.Report(log.SrcWS, wsopInbound, co.String(), eo.String(), wsocSuccess, q.ID, n, strconv.FormatInt(int64(len(msg)), 10))
				handlers.Dispatcher.Send(event)
			} else {
				eo := handlers.EventOrdinal.Next()
				log.Report(log.SrcWS, wsopInbound, co.String(), eo.String(), wsocNoRequestMapping, q.ID, n, strconv.FormatInt(int64(len(msg)), 10))
				if writeQuery(conn, co, eo, &Query{Type: queryUnexpected, ID: q.ID}) {
					break
				}
			}
		}
		close(closeCh)
	}()

Exit:
	for {
		select {
		case event := <-writeCh:
			var eo handlers.Ordinal
			var eid string
			if ti, ok := event.(handlers.TraceInfo); ok {
				eo, eid = ti.Ordinal(), ti.TraceID()
			}
			if query := queryFromEvent(event); query != nil {
				if writeQuery(conn, co, eo, query) {
					break Exit
				}
			} else {
				log.Report(log.SrcWS, wsopOutbound, co.String(), eo.String(), wsocUnexpectedResponse, eid, fmt.Sprintf("%T", event))
			}
		case <-closeCh:
			break Exit
		case <-stopCh:
			break Exit
		}
	}
	cancel()
}

func writeQuery(conn net.Conn, co handlers.Ordinal, eo handlers.Ordinal, q *Query) bool {
	w := newWebSocketTextWriter(conn)
	n := queryNameMap[q.Type]
	if err := json.NewEncoder(w).Encode(q); err != nil {
		if isConnClosed(err) {
			return true
		}
		log.Report(log.SrcWS, wsopOutbound, co.String(), eo.String(), wsocWriteError, q.ID, n, err.Error())
	} else {
		log.Report(log.SrcWS, wsopOutbound, co.String(), eo.String(), wsocSuccess, q.ID, n, strconv.FormatUint(w.length, 10))
	}
	return false
}

func isConnClosed(err error) bool {
	return err == io.EOF || err == net.ErrClosed
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
