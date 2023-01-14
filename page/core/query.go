package core

import (
	"encoding/json"
	"net/url"

	"github.com/stas-makutin/howeve/api"
)

var queryTypeToPath = map[api.QueryType]string{
	api.QueryRestart:            "/restart",
	api.QueryGetConfig:          "/cfg",
	api.QueryProtocolList:       "/protocols",
	api.QueryTransportList:      "/transports",
	api.QueryProtocolInfo:       "/protocolInfo",
	api.QueryProtocolDiscover:   "/discover",
	api.QueryProtocolDiscovery:  "/discovery",
	api.QueryAddService:         "/service/add",
	api.QueryRemoveService:      "/service/remove",
	api.QueryChangeServiceAlias: "/service/alias",
	api.QueryServiceStatus:      "/service/status",
	api.QueryListServices:       "/service/list",
	api.QuerySendToService:      "/service/send",
	api.QueryGetMessage:         "/messages/get",
	api.QueryListMessages:       "/messages/list",
}

func StringToQuery(s string) (*api.Query, error) {
	var q api.Query
	err := json.Unmarshal([]byte(s), &q)
	return &q, err
}

func QueryToString(q interface{}) (string, error) {
	b, err := json.Marshal(q)
	return string(b), err
}

func FetchQuery(url string, r interface{}, then func(r *api.Query), catch func(err string)) {
	var fi *FetchInit
	if r != nil {
		body, err := QueryToString(r)
		if err != nil {
			catch("Unexpected request format: " + err.Error())
			return
		}
		fi = &FetchInit{Method: "POST", Body: body}
	}
	Fetch(url, fi, func(response *FetchResponse) {
		var errMsg string
		if response.OK {
			q, err := StringToQuery(response.Body)
			if err == nil {
				then(q)
				return
			}
			errMsg = "Unexpected response format: " + err.Error()
		} else {
			errMsg = "Unexpected response: " + response.Body
		}
		catch(errMsg)
	}, func(err *FetchError) {
		errMsg := err.Message
		if errMsg == "" {
			errMsg = "Unable to collect the requested data"
		}
		catch(errMsg)
	})
}

func FetchQueryWithSocket(r *api.Query, then func(r *api.Query), catch func(err string)) {
	request, err := QueryToString(r)
	if err != nil {
		catch("Unexpected request format: " + err.Error())
		return
	}
	socket := NewWebSocket(WebSocketUrl())
	socket.OnOpen(func() {
		socket.Send(request)
	})
	socket.OnMessage(func(data string) {
		socket.Close()
		if q, err := StringToQuery(data); err == nil {
			then(q)
		} else {
			catch("Unexpected response format: " + err.Error())
		}
	})
	socket.OnError(func() {
		socket.Close()
		catch("Unable to collect the requested data")
	})
}

func Query(useSocket bool, r *api.Query, then func(*api.Query), catch func(err string)) {
	path, ok := queryTypeToPath[r.Type]
	if !ok {
		catch("Unexpected query")
		return
	}

	if useSocket {
		FetchQueryWithSocket(r, then, catch)
		return
	}

	u := HTTPUrl(path)
	if r.ID != "" {
		u += "?i=" + url.QueryEscape(r.ID)
	}
	FetchQuery(u, r.Payload, then, catch)
}

type CachedQuery[T any] struct {
	Value *T
	Error string
}

func (cq *CachedQuery[T]) Query(useSocket, force bool, r *api.Query, handle func(*api.Query) (*T, string), format func(string) string, success func(*T), failure func(err string)) bool {
	if force || (cq.Value == nil && cq.Error == "") {
		Query(
			useSocket, r,
			func(q *api.Query) {
				cq.Value, cq.Error = handle(q)
				if cq.Error != "" {
					cq.Error = format(cq.Error)
					failure(cq.Error)
				} else {
					success(cq.Value)
				}
			},
			func(err string) {
				cq.Error = format(err)
				failure(cq.Error)
			},
		)
		return false
	}

	if cq.Error != "" {
		failure(cq.Error)
	} else {
		success(cq.Value)
	}
	return true
}
