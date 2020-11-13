package httpsrv

import "net/http"

func setupRoutes(mux *http.ServeMux) {

	mux.Handle("/socket", handlerCtxFunc(handleWebsocket))

	mux.Handle("/restart", handlerFunc(handleRestart))
	mux.Handle("/cfg", handlerFunc(handleConfig))
}
