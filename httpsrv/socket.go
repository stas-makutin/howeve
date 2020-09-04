package httpsrv

import (
	"fmt"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		// handle error
	}
	//go messageLoop(conn)

	//go func() {
	defer conn.Close()

	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			// handle error
		}
		err = wsutil.WriteServerMessage(conn, op, msg)
		if err != nil {
			// handle error
		}
		err = wsutil.WriteServerMessage(conn, ws.OpPing, nil)
		if err != nil {
			fmt.Println(err)
		}

	}
	//}()
}

/*
func messageLoop(conn net.Conn) {
	defer conn.Close()

	ch := make(chan interface{}, 1000)
	id := eventh.Dispatcher.Subscribe(func(event interface{}) {

	})
	defer eventh.Dispatcher.Unsubscribe(id)

}
*/
