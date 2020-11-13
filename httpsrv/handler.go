package httpsrv

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/stas-makutin/howeve/events"
	"github.com/stas-makutin/howeve/events/handlers"
)

func handleEvents(w http.ResponseWriter, r *http.Request, responseType reflect.Type, request func(*http.Request) (events.TargetedRequest, bool)) {
	req, initTrace := request(r)
	if initTrace {
		if ts, ok := req.(handlers.TraceSet); ok {
			ts.InitTrace(r.URL.Query().Get("i"))
		}
	}

	var eo handlers.Ordinal
	var traceID string
	if ti, ok := req.(handlers.TraceInfo); ok {
		eo, traceID = ti.Ordinal(), ti.TraceID()
	}
	appendLogFields(r, eo.String(), traceID)

	handlers.Dispatcher.RequestResponse(r.Context(), req, responseType, func(event interface{}) {
		if query := queryFromEvent(event); query != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			err := json.NewEncoder(w).Encode(query)
			if err != nil {
				appendLogFields(r, err.Error())
			}
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	})
}

func handleRestart(w http.ResponseWriter, r *http.Request) {
	handleEvents(w, r, reflect.TypeOf(&handlers.RestartResult{}), func(h *http.Request) (events.TargetedRequest, bool) {
		return &handlers.Restart{}, true
	})
}

func handleConfig(w http.ResponseWriter, r *http.Request) {
	handleEvents(w, r, reflect.TypeOf(&handlers.ConfigGetResult{}), func(h *http.Request) (events.TargetedRequest, bool) {
		return &handlers.ConfigGet{}, true
	})
}
