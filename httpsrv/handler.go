package httpsrv

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/api"
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
	var q *api.ProtocolInfo
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &api.ProtocolInfo{}
		for _, v := range r.Form["protocols"] {
			for _, vp := range strings.FieldsFunc(v, func(c rune) bool { return c == ',' || c == ';' || c == ':' || c == '|' }) {
				n, err := strconv.ParseUint(vp, 10, 8)
				if err != nil {
					return nil, true, err
				}
				q.Protocols = append(q.Protocols, api.ProtocolIdentifier(n))
			}
		}
		for _, v := range r.Form["transports"] {
			for _, vp := range strings.FieldsFunc(v, func(c rune) bool { return c == ',' || c == ';' || c == ':' || c == '|' }) {
				n, err := strconv.ParseUint(vp, 10, 8)
				if err != nil {
					return nil, true, err
				}
				q.Transports = append(q.Transports, api.TransportIdentifier(n))
			}
		}
	}
	return &handlers.ProtocolInfo{ProtocolInfo: q}, true, nil
}

func parseProtocolDiscover(w http.ResponseWriter, r *http.Request) (events.TargetedRequest, bool, error) {
	var q *api.ProtocolDiscover
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &api.ProtocolDiscover{}

		n, err := strconv.ParseUint(r.Form.Get("protocol"), 10, 8)
		if err != nil {
			return nil, true, err
		}
		q.Protocol = api.ProtocolIdentifier(n)

		n, err = strconv.ParseUint(r.Form.Get("transport"), 10, 8)
		if err != nil {
			return nil, true, err
		}
		q.Transport = api.TransportIdentifier(n)

		if pi, ok := defs.Protocols[q.Protocol]; ok {
			if pti, ok := pi.Transports[q.Transport]; ok {
				for name, p := range pti.DiscoveryParams {
					if p.Flags&defs.ParamFlagConst == 0 {
						v := r.Form.Get(name)
						if v != "" {
							if q.Params == nil {
								q.Params = make(api.RawParamValues)
							}
							q.Params[name] = v
						}
					}
				}
			}
		}
	}

	return &handlers.ProtocolDiscover{ProtocolDiscover: q}, true, nil
}

func parseProtocolDiscovery(w http.ResponseWriter, r *http.Request) (events.TargetedRequest, bool, error) {
	var q *api.ProtocolDiscovery
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &api.ProtocolDiscovery{}

		id, err := uuid.Parse(r.Form.Get("id"))
		if err != nil {
			return nil, true, err
		}
		q.ID = id

		stop := strings.ToLower(r.Form.Get("stop"))
		q.Stop = stop == "true" || stop == "1" || stop == "yes"
	}

	return &handlers.ProtocolDiscovery{ProtocolDiscovery: q}, true, nil
}

func parseAddService(w http.ResponseWriter, r *http.Request) (events.TargetedRequest, bool, error) {
	var q *api.ServiceEntry
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &api.ServiceEntry{}

		n, err := strconv.ParseUint(r.Form.Get("protocol"), 10, 8)
		if err != nil {
			return nil, true, err
		}
		q.Protocol = api.ProtocolIdentifier(n)

		n, err = strconv.ParseUint(r.Form.Get("transport"), 10, 8)
		if err != nil {
			return nil, true, err
		}
		q.Transport = api.TransportIdentifier(n)

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
									q.Params = make(api.RawParamValues)
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
	var q *api.ServiceID
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &api.ServiceID{}

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
				q.ServiceKey = &api.ServiceKey{
					Protocol:  api.ProtocolIdentifier(pid),
					Transport: api.TransportIdentifier(tid),
					Entry:     r.Form.Get("entry"),
				}
			}
		}
		q.Alias = r.Form.Get("alias")
	}
	return &handlers.RemoveService{ServiceID: q}, true, nil
}

func parseChangeServiceAlias(w http.ResponseWriter, r *http.Request) (events.TargetedRequest, bool, error) {
	var q *api.ChangeServiceAlias
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &api.ChangeServiceAlias{}

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
				q.ServiceKey = &api.ServiceKey{
					Protocol:  api.ProtocolIdentifier(pid),
					Transport: api.TransportIdentifier(tid),
					Entry:     r.Form.Get("entry"),
				}
			}
		}
		q.Alias = r.Form.Get("alias")
		q.NewAlias = r.Form.Get("newAlias")
	}
	return &handlers.ChangeServiceAlias{ChangeServiceAlias: q}, true, nil
}

func parseServiceStatus(w http.ResponseWriter, r *http.Request) (events.TargetedRequest, bool, error) {
	var q *api.ServiceID
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &api.ServiceID{}

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
				q.ServiceKey = &api.ServiceKey{
					Protocol:  api.ProtocolIdentifier(pid),
					Transport: api.TransportIdentifier(tid),
					Entry:     r.Form.Get("entry"),
				}
			}
		}
		q.Alias = r.Form.Get("alias")
	}
	return &handlers.ServiceStatus{ServiceID: q}, true, nil
}

func parseListServices(w http.ResponseWriter, r *http.Request) (events.TargetedRequest, bool, error) {
	var q *api.ListServices
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &api.ListServices{}
		for _, v := range r.Form["protocols"] {
			for _, vp := range strings.FieldsFunc(v, func(c rune) bool { return c == ',' || c == ';' || c == ':' || c == '|' }) {
				n, err := strconv.ParseUint(vp, 10, 8)
				if err != nil {
					return nil, true, err
				}
				q.Protocols = append(q.Protocols, api.ProtocolIdentifier(n))
			}
		}
		for _, v := range r.Form["transports"] {
			for _, vp := range strings.FieldsFunc(v, func(c rune) bool { return c == ',' || c == ';' || c == ':' || c == '|' }) {
				n, err := strconv.ParseUint(vp, 10, 8)
				if err != nil {
					return nil, true, err
				}
				q.Transports = append(q.Transports, api.TransportIdentifier(n))
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
	return &handlers.ListServices{ListServices: q}, true, nil
}

func parseSendToService(w http.ResponseWriter, r *http.Request) (events.TargetedRequest, bool, error) {
	var q *api.SendToService
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &api.SendToService{}

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
				q.ServiceKey = &api.ServiceKey{
					Protocol:  api.ProtocolIdentifier(pid),
					Transport: api.TransportIdentifier(tid),
					Entry:     r.Form.Get("entry"),
				}
			}
		}
		q.Alias = r.Form.Get("alias")
		if q.Payload, err = base64.StdEncoding.DecodeString(r.Form.Get("payload")); err != nil {
			return nil, true, err
		}
	}
	return &handlers.SendToService{SendToService: q}, true, nil
}

func parseGetMessage(w http.ResponseWriter, r *http.Request) (events.TargetedRequest, bool, error) {
	var q uuid.UUID
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q, err = uuid.Parse(r.Form.Get("id"))
		if err != nil {
			return nil, true, err
		}
	}
	return &handlers.GetMessage{ID: q}, true, nil
}

func parseListMessages(w http.ResponseWriter, r *http.Request) (events.TargetedRequest, bool, error) {
	var q *api.ListMessages
	if ok, err := parseJSONRequest(&q, w, r, 4096); ok {
		if err != nil {
			return nil, true, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, true, err
		}
		q = &api.ListMessages{}
		n, err := strconv.ParseInt(r.Form.Get("index"), 10, 8)
		if err != nil {
			return nil, true, err
		}
		index := int(n)
		q.FromIndex = &index

		n, err = strconv.ParseInt(r.Form.Get("count"), 10, 8)
		if err != nil {
			return nil, true, err
		}
		q.Count = int(n)
	}
	return &handlers.ListMessages{ListMessages: q}, true, nil
}
