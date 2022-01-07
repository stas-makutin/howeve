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
					return parseProtocolInfo(w, r)
				})
			},
		},
		{
			"/discover", func(w http.ResponseWriter, r *http.Request) {
				handleEvents(w, r, reflect.TypeOf(&handlers.ProtocolDiscoverResult{}), func(h *http.Request) (events.TargetedRequest, bool, error) {
					return parseProtocolDiscover(w, r)
				})
			},
		},
		{
			"/discovery", func(w http.ResponseWriter, r *http.Request) {
				handleEvents(w, r, reflect.TypeOf(&handlers.ProtocolDiscoveryResult{}), func(h *http.Request) (events.TargetedRequest, bool, error) {
					return parseProtocolDiscovery(w, r)
				})
			},
		},
		{
			"/service/add", func(w http.ResponseWriter, r *http.Request) {
				handleEvents(w, r, reflect.TypeOf(&handlers.AddServiceResult{}), func(h *http.Request) (events.TargetedRequest, bool, error) {
					return parseAddService(w, r)
				})
			},
		},
		{
			"/service/remove", func(w http.ResponseWriter, r *http.Request) {
				handleEvents(w, r, reflect.TypeOf(&handlers.RemoveServiceResult{}), func(h *http.Request) (events.TargetedRequest, bool, error) {
					return parseRemoveService(w, r)
				})
			},
		},
		{
			"/service/alias", func(w http.ResponseWriter, r *http.Request) {
				handleEvents(w, r, reflect.TypeOf(&handlers.ChangeServiceAliasResult{}), func(h *http.Request) (events.TargetedRequest, bool, error) {
					return parseChangeServiceAlias(w, r)
				})
			},
		},
		{
			"/service/status", func(w http.ResponseWriter, r *http.Request) {
				handleEvents(w, r, reflect.TypeOf(&handlers.ServiceStatusResult{}), func(h *http.Request) (events.TargetedRequest, bool, error) {
					return parseServiceStatus(w, r)
				})
			},
		},
		{
			"/service/list", func(w http.ResponseWriter, r *http.Request) {
				handleEvents(w, r, reflect.TypeOf(&handlers.ListServicesResult{}), func(h *http.Request) (events.TargetedRequest, bool, error) {
					return parseListServices(w, r)
				})
			},
		},
		{
			"/service/send", func(w http.ResponseWriter, r *http.Request) {
				handleEvents(w, r, reflect.TypeOf(&handlers.SendToServiceResult{}), func(h *http.Request) (events.TargetedRequest, bool, error) {
					return parseSendToService(w, r)
				})
			},
		},
		{
			"/messages/get", func(w http.ResponseWriter, r *http.Request) {
				handleEvents(w, r, reflect.TypeOf(&handlers.GetMessageResult{}), func(h *http.Request) (events.TargetedRequest, bool, error) {
					return parseGetMessage(w, r)
				})
			},
		},
		{
			"/messages/info", func(w http.ResponseWriter, r *http.Request) {
				handleEvents(w, r, reflect.TypeOf(&handlers.GetMessagesInfoResult{}), func(h *http.Request) (events.TargetedRequest, bool, error) {
					return &handlers.GetMessagesInfo{}, true, nil
				})
			},
		},
		{
			"/messages/after", func(w http.ResponseWriter, r *http.Request) {
				handleEvents(w, r, reflect.TypeOf(&handlers.MessagesAfterResult{}), func(h *http.Request) (events.TargetedRequest, bool, error) {
					return parseMessagesAfter(w, r)
				})
			},
		},
		{
			"/messages/list", func(w http.ResponseWriter, r *http.Request) {
				handleEvents(w, r, reflect.TypeOf(&handlers.ListMessagesResult{}), func(h *http.Request) (events.TargetedRequest, bool, error) {
					return parseListMessages(w, r)
				})
			},
		},
	} {
		mux.Handle(rt.route, handlerFunc(rt.handler))
	}
}
