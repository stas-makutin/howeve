package httpsrv

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/stas-makutin/howeve/defs"
	"github.com/stas-makutin/howeve/events"
	"github.com/stas-makutin/howeve/events/handlers"
	"github.com/stas-makutin/howeve/services"
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

func parseJSONRequest(v interface{}, w http.ResponseWriter, r *http.Request, maxSize int64) (bool, error) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		return false, nil
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)
	d := json.NewDecoder(r.Body)
	if err := d.Decode(v); err != nil {
		return true, err
	}
	return true, nil
}

func parseProtocolInfo(w http.ResponseWriter, r *http.Request) (events.TargetedRequest, bool, error) {
	var q *handlers.ProtocolInfoFilter
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &handlers.ProtocolInfoFilter{}
		for _, v := range r.Form["protocols"] {
			for _, vp := range strings.FieldsFunc(v, func(c rune) bool { return c == ',' || c == ';' || c == ':' || c == '|' }) {
				if n, err := strconv.ParseUint(vp, 10, 8); err != nil {
					return nil, true, err
				} else {
					q.Protocols = append(q.Protocols, defs.ProtocolIdentifier(n))
				}
			}
		}
		for _, v := range r.Form["transports"] {
			for _, vp := range strings.FieldsFunc(v, func(c rune) bool { return c == ',' || c == ';' || c == ':' || c == '|' }) {
				if n, err := strconv.ParseUint(vp, 10, 8); err != nil {
					return nil, true, err
				} else {
					q.Transports = append(q.Transports, defs.TransportIdentifier(n))
				}
			}
		}
	}
	return &handlers.ProtocolInfo{Filter: q}, true, nil
}

func parseProtocolDiscovery(w http.ResponseWriter, r *http.Request) (events.TargetedRequest, bool, error) {
	var q *handlers.ProtocolDiscoveryQuery
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &handlers.ProtocolDiscoveryQuery{}
		if n, err := strconv.ParseUint(r.Form.Get("protocol"), 10, 8); err != nil {
			return nil, true, err
		} else {
			q.Protocol = defs.ProtocolIdentifier(n)
		}
		if n, err := strconv.ParseUint(r.Form.Get("transport"), 10, 8); err != nil {
			return nil, true, err
		} else {
			q.Transport = defs.TransportIdentifier(n)
		}
		if pi, ok := services.Protocols[q.Protocol]; ok {
			if pti, ok := pi.Transports[q.Transport]; ok {
				for name, p := range pti.DiscoveryParams {
					if p.Flags&defs.ParamFlagConst == 0 {
						v := r.Form.Get(name)
						if v != "" {
							if q.Params == nil {
								q.Params = make(handlers.ParamsValues)
							}
							q.Params[name] = v
						}
					}
				}
			}
		}
	}

	return &handlers.ProtocolDiscovery{ProtocolDiscoveryQuery: q}, true, nil
}

func parseAddService(w http.ResponseWriter, r *http.Request) (events.TargetedRequest, bool, error) {
	var q *handlers.ServiceEntry
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &handlers.ServiceEntry{}
		if n, err := strconv.ParseUint(r.Form.Get("protocol"), 10, 8); err != nil {
			return nil, true, err
		} else {
			q.Protocol = defs.ProtocolIdentifier(n)
		}
		if n, err := strconv.ParseUint(r.Form.Get("transport"), 10, 8); err != nil {
			return nil, true, err
		} else {
			q.Transport = defs.TransportIdentifier(n)
		}
		q.Entry = r.Form.Get("entry")
		if q.Entry == "" {
			return nil, true, errors.New("'entry' parameter is not available or empty")
		}
		if pi, ok := services.Protocols[q.Protocol]; ok {
			if pti, ok := pi.Transports[q.Transport]; ok {
				if ti, ok := services.Transports[q.Transport]; ok {
					for name, p := range pti.Params.Merge(ti.Params) {
						if p.Flags&defs.ParamFlagConst == 0 {
							v := r.Form.Get(name)
							if v != "" {
								if q.Params == nil {
									q.Params = make(handlers.ParamsValues)
								}
								q.Params[name] = v
							}
						}
					}
				}
			}
		}
	}
	return &handlers.AddService{ServiceEntry: q}, true, nil
}
