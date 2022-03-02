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
}

type ProtocolsLoadFailed struct {
}

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
}

func protocolsLoad(action *ProtocolsLoad) bool {
	if action.Force || pvStore.Data == false {
		if action.UseSocket {
			protocolsLoadSockets()
		} else {
			protocolsLoadFetch()
		}
		// return false TODO - return false on async op
	}
	core.Dispatch(&ProtocolsLoaded{})
	return true
}

// store

type ProtocolViewStore struct {
	Loading   bool
	UseSocket bool
	Data      bool
}

var pvStore = &ProtocolViewStore{
	Loading:   true,
	UseSocket: true,
	Data:      false,
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
		if protocolsLoad(e) {
			return
		}
	case *ProtocolsLoaded:
		pvStore.Loading = false
		pvStore.Data = true
	case *ProtocolsLoadFailed:
		pvStore.Loading = false
		// TODO
	default:
		return
	}
	core.Dispatch(ChangeEvent{pvStore})
}
