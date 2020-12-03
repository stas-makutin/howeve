package httpsrv

import (
	"net/http"
	"reflect"

	"github.com/stas-makutin/howeve/events"
	"github.com/stas-makutin/howeve/events/handlers"
)

func setupRoutes(mux *http.ServeMux) {

	mux.Handle("/socket", handlerCtxFunc(handleWebsocket))

	for _, rt := range []struct {
		route   string
		handler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			"/restart", func(w http.ResponseWriter, r *http.Request) {
				handleEvents(w, r, reflect.TypeOf(&handlers.RestartResult{}), func(h *http.Request) (events.TargetedRequest, bool, error) {
					return &handlers.Restart{}, true, nil
				})
			},
		},
		{
			"/cfg", func(w http.ResponseWriter, r *http.Request) {
				handleEvents(w, r, reflect.TypeOf(&handlers.ConfigGetResult{}), func(h *http.Request) (events.TargetedRequest, bool, error) {
					return &handlers.ConfigGet{}, true, nil
				})
			},
		},
		{
			"/protocols", func(w http.ResponseWriter, r *http.Request) {
				handleEvents(w, r, reflect.TypeOf(&handlers.ProtocolListResult{}), func(h *http.Request) (events.TargetedRequest, bool, error) {
					return &handlers.ProtocolList{}, true, nil
				})
			},
		},
		{
			"/transports", func(w http.ResponseWriter, r *http.Request) {
				handleEvents(w, r, reflect.TypeOf(&handlers.TransportListResult{}), func(h *http.Request) (events.TargetedRequest, bool, error) {
					return &handlers.TransportList{}, true, nil
				})
			},
		},
		{
			"/protocolInfo", func(w http.ResponseWriter, r *http.Request) {
				handleEvents(w, r, reflect.TypeOf(&handlers.ProtocolInfoResult{}), func(h *http.Request) (events.TargetedRequest, bool, error) {
					return parseProtocolInfo(r)
				})
			},
		},
	} {
		mux.Handle(rt.route, handlerFunc(rt.handler))
	}
}
