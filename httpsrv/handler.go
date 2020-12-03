package httpsrv

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/stas-makutin/howeve/events"
	"github.com/stas-makutin/howeve/events/handlers"
)

func handleEvents(w http.ResponseWriter, r *http.Request, responseType reflect.Type, request func(*http.Request) (events.TargetedRequest, bool, error)) {
	req, initTrace, err := request(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
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

func parseProtocolInfo(r *http.Request) (events.TargetedRequest, bool, error) {
	if err := r.ParseForm(); err != nil {
		return nil, true, err
	}
	tr := &handlers.ProtocolInfo{Filter: &handlers.ProtocolInfoFilter{}}
	for _, v := range r.Form["protocols"] {
		for _, vp := range strings.FieldsFunc(v, func(c rune) bool { return c == ',' || c == ';' || c == ':' || c == '|' }) {
			if n, err := strconv.ParseUint(vp, 10, 8); err != nil {
				return nil, true, err
			} else {
				tr.Filter.Protocols = append(tr.Filter.Protocols, uint8(n))
			}
		}
	}
	for _, v := range r.Form["transports"] {
		for _, vp := range strings.FieldsFunc(v, func(c rune) bool { return c == ',' || c == ';' || c == ':' || c == '|' }) {
			if n, err := strconv.ParseUint(vp, 10, 8); err != nil {
				return nil, true, err
			} else {
				tr.Filter.Transports = append(tr.Filter.Transports, uint8(n))
			}
		}
	}
	return tr, true, nil
}
