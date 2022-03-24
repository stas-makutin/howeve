package actions

import (
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

func protocolsProcessResponse(r *api.Query) {
	if r.Payload == nil {
		core.Dispatch(&ProtocolsLoaded{&api.ProtocolInfoResult{}})
	} else if p, ok := r.Payload.(*api.ProtocolInfoResult); ok {
		core.Dispatch(&ProtocolsLoaded{p})
	} else {
		core.Dispatch(ProtocolsLoadFailed("Unexpected response type"))
	}
}

func protocolsLoadSockets() {
	core.FetchQueryWithSocket(
		&api.Query{Type: api.QueryProtocolInfo},
		func(r *api.Query) {
			protocolsProcessResponse(r)
		},
		func(err string) {
			core.Dispatch(ProtocolsLoadFailed(err))
		},
	)
}

func protocolsLoadFetch() {
	core.FetchQuery(
		core.HTTPUrl("/protocolInfo"), nil,
		func(r *api.Query) {
			protocolsProcessResponse(r)
		},
		func(err string) {
			core.Dispatch(ProtocolsLoadFailed(err))
		},
	)
}

func protocolsLoad(action *ProtocolsLoad) bool {
	if action.Force || (pvStore.Protocols == nil && pvStore.Error == "") {
		if action.UseSocket {
			protocolsLoadSockets()
		} else {
			protocolsLoadFetch()
		}
		return false
	}
	// restore saved state
	if pvStore.Error != "" {
		core.Dispatch(ProtocolsLoadFailed(pvStore.Error))
	} else {
		core.Dispatch(&ProtocolsLoaded{pvStore.Protocols})
	}
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
		pvStore.Error = ""
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
