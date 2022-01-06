package httpsrv

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/defs"
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
				n, err := strconv.ParseUint(vp, 10, 8)
				if err != nil {
					return nil, true, err
				}
				q.Protocols = append(q.Protocols, defs.ProtocolIdentifier(n))
			}
		}
		for _, v := range r.Form["transports"] {
			for _, vp := range strings.FieldsFunc(v, func(c rune) bool { return c == ',' || c == ';' || c == ':' || c == '|' }) {
				n, err := strconv.ParseUint(vp, 10, 8)
				if err != nil {
					return nil, true, err
				}
				q.Transports = append(q.Transports, defs.TransportIdentifier(n))
			}
		}
	}
	return &handlers.ProtocolInfo{Filter: q}, true, nil
}

func parseProtocolDiscover(w http.ResponseWriter, r *http.Request) (events.TargetedRequest, bool, error) {
	var q *handlers.ProtocolDiscoverInput
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &handlers.ProtocolDiscoverInput{}

		n, err := strconv.ParseUint(r.Form.Get("protocol"), 10, 8)
		if err != nil {
			return nil, true, err
		}
		q.Protocol = defs.ProtocolIdentifier(n)

		n, err = strconv.ParseUint(r.Form.Get("transport"), 10, 8)
		if err != nil {
			return nil, true, err
		}
		q.Transport = defs.TransportIdentifier(n)

		if pi, ok := defs.Protocols[q.Protocol]; ok {
			if pti, ok := pi.Transports[q.Transport]; ok {
				for name, p := range pti.DiscoveryParams {
					if p.Flags&defs.ParamFlagConst == 0 {
						v := r.Form.Get(name)
						if v != "" {
							if q.Params == nil {
								q.Params = make(defs.RawParamValues)
							}
							q.Params[name] = v
						}
					}
				}
			}
		}
	}

	return &handlers.ProtocolDiscover{ProtocolDiscoverInput: q}, true, nil
}

func parseProtocolDiscovery(w http.ResponseWriter, r *http.Request) (events.TargetedRequest, bool, error) {
	var q *handlers.ProtocolDiscoveryInput
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &handlers.ProtocolDiscoveryInput{}

		id, err := uuid.Parse(r.Form.Get("id"))
		if err != nil {
			return nil, true, err
		}
		q.ID = id

		stop := strings.ToLower(r.Form.Get("stop"))
		q.Stop = stop == "true" || stop == "1" || stop == "yes"
	}

	return &handlers.ProtocolDiscovery{ProtocolDiscoveryInput: q}, true, nil
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

		n, err := strconv.ParseUint(r.Form.Get("protocol"), 10, 8)
		if err != nil {
			return nil, true, err
		}
		q.Protocol = defs.ProtocolIdentifier(n)

		n, err = strconv.ParseUint(r.Form.Get("transport"), 10, 8)
		if err != nil {
			return nil, true, err
		}
		q.Transport = defs.TransportIdentifier(n)

		q.Entry = r.Form.Get("entry")
		q.Alias = r.Form.Get("alias")
		if pi, ok := defs.Protocols[q.Protocol]; ok {
			if pti, ok := pi.Transports[q.Transport]; ok {
				if ti, ok := defs.Transports[q.Transport]; ok {
					for name, p := range pti.Params.Merge(ti.Params) {
						if p.Flags&defs.ParamFlagConst == 0 {
							v := r.Form.Get(name)
							if v != "" {
								if q.Params == nil {
									q.Params = make(defs.RawParamValues)
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

func parseRemoveService(w http.ResponseWriter, r *http.Request) (events.TargetedRequest, bool, error) {
	var q *handlers.ServiceID
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &handlers.ServiceID{}

		if protocol := r.Form.Get("protocol"); protocol != "" {
			if transport := r.Form.Get("transport"); transport != "" {
				pid, err := strconv.ParseUint(protocol, 10, 8)
				if err != nil {
					return nil, true, err
				}
				tid, err := strconv.ParseUint(transport, 10, 8)
				if err != nil {
					return nil, true, err
				}
				q.ServiceKey = &defs.ServiceKey{
					Protocol:  defs.ProtocolIdentifier(pid),
					Transport: defs.TransportIdentifier(tid),
					Entry:     r.Form.Get("entry"),
				}
			}
		}
		q.Alias = r.Form.Get("alias")
	}
	return &handlers.RemoveService{ServiceID: q}, true, nil
}

func parseChangeServiceAlias(w http.ResponseWriter, r *http.Request) (events.TargetedRequest, bool, error) {
	var q *handlers.ChangeServiceAliasQuery
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &handlers.ChangeServiceAliasQuery{}

		if protocol := r.Form.Get("protocol"); protocol != "" {
			if transport := r.Form.Get("transport"); transport != "" {
				pid, err := strconv.ParseUint(protocol, 10, 8)
				if err != nil {
					return nil, true, err
				}
				tid, err := strconv.ParseUint(transport, 10, 8)
				if err != nil {
					return nil, true, err
				}
				q.ServiceKey = &defs.ServiceKey{
					Protocol:  defs.ProtocolIdentifier(pid),
					Transport: defs.TransportIdentifier(tid),
					Entry:     r.Form.Get("entry"),
				}
			}
		}
		q.Alias = r.Form.Get("alias")
		q.NewAlias = r.Form.Get("newAlias")
	}
	return &handlers.ChangeServiceAlias{ChangeServiceAliasQuery: q}, true, nil
}

func parseServiceStatus(w http.ResponseWriter, r *http.Request) (events.TargetedRequest, bool, error) {
	var q *handlers.ServiceID
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &handlers.ServiceID{}

		if protocol := r.Form.Get("protocol"); protocol != "" {
			if transport := r.Form.Get("transport"); transport != "" {
				pid, err := strconv.ParseUint(protocol, 10, 8)
				if err != nil {
					return nil, true, err
				}
				tid, err := strconv.ParseUint(transport, 10, 8)
				if err != nil {
					return nil, true, err
				}
				q.ServiceKey = &defs.ServiceKey{
					Protocol:  defs.ProtocolIdentifier(pid),
					Transport: defs.TransportIdentifier(tid),
					Entry:     r.Form.Get("entry"),
				}
			}
		}
		q.Alias = r.Form.Get("alias")
	}
	return &handlers.ServiceStatus{ServiceID: q}, true, nil
}

func parseListServices(w http.ResponseWriter, r *http.Request) (events.TargetedRequest, bool, error) {
	var q *handlers.ListServicesInput
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &handlers.ListServicesInput{}
		for _, v := range r.Form["protocols"] {
			for _, vp := range strings.FieldsFunc(v, func(c rune) bool { return c == ',' || c == ';' || c == ':' || c == '|' }) {
				n, err := strconv.ParseUint(vp, 10, 8)
				if err != nil {
					return nil, true, err
				}
				q.Protocols = append(q.Protocols, defs.ProtocolIdentifier(n))
			}
		}
		for _, v := range r.Form["transports"] {
			for _, vp := range strings.FieldsFunc(v, func(c rune) bool { return c == ',' || c == ';' || c == ':' || c == '|' }) {
				n, err := strconv.ParseUint(vp, 10, 8)
				if err != nil {
					return nil, true, err
				}
				q.Transports = append(q.Transports, defs.TransportIdentifier(n))
			}
		}
		for _, v := range r.Form["entries"] {
			for _, vp := range strings.FieldsFunc(v, func(c rune) bool { return c == ',' || c == ';' || c == ':' || c == '|' }) {
				if vp != "" {
					q.Entries = append(q.Entries, vp)
				}
			}
		}
		for _, v := range r.Form["aliases"] {
			for _, vp := range strings.FieldsFunc(v, func(c rune) bool { return c == ',' || c == ';' || c == ':' || c == '|' }) {
				if vp != "" {
					q.Aliases = append(q.Aliases, vp)
				}
			}
		}
	}
	return &handlers.ListServices{ListServicesInput: q}, true, nil
}

func parseSendToService(w http.ResponseWriter, r *http.Request) (events.TargetedRequest, bool, error) {
	var q *handlers.SendToServiceInput
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &handlers.SendToServiceInput{}

		if protocol := r.Form.Get("protocol"); protocol != "" {
			if transport := r.Form.Get("transport"); transport != "" {
				pid, err := strconv.ParseUint(protocol, 10, 8)
				if err != nil {
					return nil, true, err
				}
				tid, err := strconv.ParseUint(transport, 10, 8)
				if err != nil {
					return nil, true, err
				}
				q.ServiceKey = &defs.ServiceKey{
					Protocol:  defs.ProtocolIdentifier(pid),
					Transport: defs.TransportIdentifier(tid),
					Entry:     r.Form.Get("entry"),
				}
			}
		}
		q.Alias = r.Form.Get("alias")
		if q.Payload, err = base64.StdEncoding.DecodeString(r.Form.Get("payload")); err != nil {
			return nil, true, err
		}
	}
	return &handlers.SendToService{SendToServiceInput: q}, true, nil
}
