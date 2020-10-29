package httpsrv

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/stas-makutin/howeve/events"
	"github.com/stas-makutin/howeve/events/handlers"
)

func handleWebsocket(w http.ResponseWriter, r *http.Request, hc *handlerContext) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		appendLogFields(r, fmt.Sprintf("%T: %v", err, err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	co := nextWsConnOrdinal()
	co.logOpen()
	closeCh := make(chan struct{})
	go func() {
		select {
		case <-hc.stopCh: // closes connection on http server shutdown
		case <-closeCh: // closes connection if it is terminated
		}
		conn.Close()
		co.logClose()
	}()
	hc.handlerWg.Add(1)
	go func() {
		defer hc.handlerWg.Done()
		defer close(closeCh)
		messageLoop(conn, co)
	}()
}

func messageLoop(conn net.Conn, co wsConnOrdinal) {
	var id events.SubscriberID
	id = handlers.Dispatcher.Subscribe(func(event interface{}) {
		if te, ok := event.(events.TargetedResponse); ok && te.Receiver() == id {
			switch te.(type) {
			case *handlers.ConfigData:
				co.logMsg(0, "", false, string(queryGetConfig), 0)

				err := json.NewEncoder(WebSocketTextWriter{conn}).Encode(event.(*handlers.ConfigData).Config)
				if err != nil {
					// TODO
				}
			}
		}
	})
	defer handlers.Dispatcher.Unsubscribe(id)

	for {
		mo := nextWsMsgOrdinal()

		msg, _, err := wsutil.ReadClientData(conn)
		if err != nil {
			if err != io.EOF {
				// TODO report error
			}
			break
		}

		var q WebSocketQuery
		err = json.Unmarshal(msg, &q)
		if err != nil {
			// TODO report error
			continue
		}

		co.logMsg(mo, q.ID, true, string(q.Type), int64(len(msg)))

		// process message
		switch q.Type {
		case queryGetConfig:
			handlers.Dispatcher.Send(&handlers.ConfigGet{RequestTarget: events.RequestTarget{ReceiverID: id}})
		}
	}
}

// WebSocketTextWriter struct
type WebSocketTextWriter struct {
	conn net.Conn
}

func (w WebSocketTextWriter) Write(p []byte) (n int, err error) {
	n = 0
	err = wsutil.WriteServerText(w.conn, p)
	return
}
