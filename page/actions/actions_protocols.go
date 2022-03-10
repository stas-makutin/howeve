package actions

import (
	"encoding/json"

	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/page/core"
)

func init() {
	core.DispatcherSubscribe(pvAction)
}

// actions

type ProtocolsUseSocket bool

type ProtocolsLoad struct {
	Force     bool
	UseSocket bool
}

type ProtocolsLoaded struct {
	Protocols *api.ProtocolInfoResult
}

type ProtocolsLoadFailed string

func protocolsLoadSockets() {
	socket := core.NewWebSocket(core.WebSocketUrl())
	socket.OnOpen(func() {
		q := &api.Query{Type: api.QueryProtocolInfo}
		b, _ := json.Marshal(q)
		socket.Send(string(b))
	})
	socket.OnMessage(func(data string) {
		core.Console.Log(data)
		socket.Close()
	})
}

func protocolsLoadFetch() {
	core.Fetch(core.HTTPUrl("/protocolInfo"), nil,
		func(response *core.FetchResponse) {
			core.Console.Log("fetch data: " + response.Body)
		},
		func(err *core.FetchError) {
			core.Console.Log("fetch error: " + err.Name + ": " + err.Message)
		},
	)
}

func protocolsLoad(action *ProtocolsLoad) bool {
	if action.Force || pvStore.Protocols == nil {
		if action.UseSocket {
			protocolsLoadSockets()
		} else {
			protocolsLoadFetch()
		}
		// return false TODO
	}
	core.Dispatch(&ProtocolsLoaded{})
	return true
}

// store

type ProtocolViewStore struct {
	Loading   bool
	UseSocket bool
	Error     string
	Protocols *api.ProtocolInfoResult
}

var pvStore = &ProtocolViewStore{
	Loading:   true,
	UseSocket: true,
}
var pvStoreChanging = false

func GetProtocolViewStore() *ProtocolViewStore {
	return pvStore
}

// reducer

func pvAction(event interface{}) {
	switch e := event.(type) {
	case ProtocolsUseSocket:
		pvStore.UseSocket = bool(e)
	case *ProtocolsLoad:
		pvStore.Loading = true
		pvStore.Error = ""
		if protocolsLoad(e) {
			return
		}
	case *ProtocolsLoaded:
		pvStore.Loading = false
		pvStore.Protocols = e.Protocols
	case ProtocolsLoadFailed:
		pvStore.Loading = false
		pvStore.Error = string(e)
		if pvStore.Error == "" {
			pvStore.Error = "Unable to load protocol information"
		}
	default:
		return
	}
	core.Dispatch(ChangeEvent{pvStore})
}
