package httpsrv

import "net/http"

func setupRoutes(mux *http.ServeMux) {

	mux.Handle("/socket", handlerCtxFunc(handleWebsocket))

	mux.Handle("/cfg", handlerFunc(handleConfig))
}
