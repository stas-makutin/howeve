package httpsrv

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/stas-makutin/howeve/eventh"
	"github.com/stas-makutin/howeve/events"
)

type queryIDType uint16

const (
	queryGetConfig = queryIDType(iota)
)

// WebSocketQuery struct
type WebSocketQuery struct {
	QueryID queryIDType `json:"q"`
	Payload interface{} `json:"p,omitempty"`
}

// UnmarshalJSON func
func (c *WebSocketQuery) UnmarshalJSON(data []byte) error {
	var env struct {
		QueryID string          `json:"q"`
		Payload json.RawMessage `json:"p,omitempty"`
	}
	err := json.Unmarshal(data, &env)
	if err != nil {
		return err
	}
	id, ok := map[string]queryIDType{
		"getCfg": queryGetConfig, "getConfig": queryGetConfig,
	}[env.QueryID]
	if !ok {
		return fmt.Errorf("Unknown query %v", env.QueryID)
	}
	c.QueryID = id
	c.Payload = nil
	/*
		switch c.QueryID {
		}
	*/
	return err
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

func handleWebsocket(w http.ResponseWriter, r *http.Request, hc *handlerContext) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		// todo handle error
	}
	go func() {
		// closes connection on http server shutdown
		hc.stopWg.Wait()
		conn.Close()
	}()
	hc.handlerWg.Add(1)
	go messageLoop(conn, hc)
}

func messageLoop(conn net.Conn, hc *handlerContext) {
	defer hc.handlerWg.Done()
	defer conn.Close()

	var id events.SubscriberID
	id = eventh.Dispatcher.Subscribe(func(event interface{}) {
		if te, ok := event.(events.TargetedResponse); ok && te.Receiver() == id {
			switch te.(type) {
			case *eventh.ConfigData:
				err := json.NewEncoder(WebSocketTextWriter{conn}).Encode(event.(*eventh.ConfigData).Config)
				if err != nil {
					// TODO
				}
			}
		}
	})
	defer eventh.Dispatcher.Unsubscribe(id)

	for {
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

		// process message
		switch q.QueryID {
		case queryGetConfig:
			eventh.Dispatcher.Send(&eventh.ConfigGet{RequestTarget: events.RequestTarget{ReceiverID: id}})
		}
	}

	fmt.Println("-- disconnected")
}
