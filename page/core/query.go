package core

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	"github.com/stas-makutin/howeve/api"
)

const (
	querySocketConnectTimeout = 3000
	querySocketSendTimeout    = 3000
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
	timeout := &Timeout{}
	timeout.Set(func() {
		socket.Close()
		timeout.Clear()
		catch("The request to the server timed out")
	}, querySocketConnectTimeout)
	socket.OnOpen(func() {
		timeout.Reset(querySocketSendTimeout)
		socket.Send(request)
	})
	socket.OnMessage(func(data string) {
		timeout.Clear()
		socket.Close()
		if q, err := StringToQuery(data); err == nil {
			then(q)
		} else {
			catch("Unexpected response format: " + err.Error())
		}
	})
	socket.OnError(func() {
		timeout.Clear()
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

func (cq *CachedQuery[T]) Query(useSocket, force bool, r *api.Query, handle func(*api.Query) (*T, string), success func(*T), failure func(err string)) bool {
	if force || (cq.Value == nil && cq.Error == "") {
		Query(
			useSocket, r,
			func(q *api.Query) {
				cq.Value, cq.Error = handle(q)
				if cq.Error != "" {
					failure(cq.Error)
				} else {
					success(cq.Value)
				}
			},
			func(err string) {
				cq.Error = err
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

var mainSocket *SharedSocket

func MainSocket() *SharedSocket {
	if mainSocket == nil {
		mainSocket = newSharedSocket()
	}
	return mainSocket
}

type MainSocketOpened struct{}
type MainSocketMessage *api.Query
type MainSocketError string

type SharedSocket struct {
	socket      *WebSocket
	openTimeout *Timeout
	ready       bool
	lastID      int
}

func newSharedSocket() *SharedSocket {
	s := &SharedSocket{}
	s.openTimeout = &Timeout{}
	s.init()
	return s
}

func (s *SharedSocket) init() {
	if s.socket != nil {
		s.close()
	}
	s.socket = NewWebSocket(WebSocketUrl())
	s.openTimeout.Set(func() {
		Dispatch(MainSocketError("The communication with server timed out"))
		s.init()
	}, querySocketConnectTimeout)

	s.socket.OnOpen(s.onOpen)
	s.socket.OnClose(s.onClose)
	s.socket.OnError(s.onError)
	s.socket.OnMessage(s.onMessage)
}

func (s *SharedSocket) Send(r *api.Query) (string, bool) {
	success := false
	id := r.ID
	if id == "" {
		s.lastID += 1
		id = strconv.Itoa(s.lastID)
		r.ID = id
	}
	if s.ready {
		request, err := QueryToString(r)
		if err == nil {
			s.socket.Send(request)
			success = true
		} else {
			Dispatch(MainSocketError("Unexpected request format: " + err.Error()))
		}
	} else {
		Dispatch(MainSocketError("The communication with server not established"))
	}
	return id, success
}

func (s *SharedSocket) SendWithTimeout(r *api.Query, t *Timeout) (string, bool) {
	t.Set(func() {
		Dispatch(MainSocketError("The communication with server timed out"))
		time.Sleep(500)
		s.init()
	}, querySocketSendTimeout)
	id, success := s.Send(r)
	if !success {
		t.Clear()
	}
	return id, success
}

func (s *SharedSocket) close() {
	s.openTimeout.Clear()
	s.socket.Close()
	s.ready = false
	s.lastID = 0

}

func (s *SharedSocket) onOpen() {
	s.openTimeout.Clear()
	s.ready = true
	Dispatch(MainSocketOpened{})
}

func (s *SharedSocket) onClose() {
	Dispatch(MainSocketError("The communication with server closed"))
	time.Sleep(500)
	s.init()
}

func (s *SharedSocket) onError() {
	Dispatch(MainSocketError("The communication with server failed"))
	time.Sleep(500)
	s.init()
}

func (s *SharedSocket) onMessage(data string) {
	s.openTimeout.Clear()
	if q, err := StringToQuery(data); err == nil {
		Dispatch(MainSocketMessage(q))
	} else {
		Dispatch(MainSocketError("Unexpected response format: " + err.Error()))
	}
}
