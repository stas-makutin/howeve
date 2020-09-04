package httpsrv

import "net/http"

func setupRoutes(mux *http.ServeMux) {

	mux.Handle("/socket", http.HandlerFunc(handleWebsocket))

	mux.Handle("/cfg", http.HandlerFunc(handleConfig))
}
