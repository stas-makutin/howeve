package actions

import (
	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/page/core"
)

func init() {
	core.DispatcherSubscribe(pvAction)
}

// store

type ProtocolViewStore struct {
	Loading      bool
	UseSocket    bool
	Protocols    core.CachedQuery[api.ProtocolInfoResult]
	DisplayError string
}

var pvStore = &ProtocolViewStore{
	Loading:   true,
	UseSocket: true,
}

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
		pvStore.DisplayError = ""
		if protocolsLoad(e.Force, e.UseSocket) {
			return
		}
	case ProtocolsLoaded:
		pvStore.Loading = false
		pvStore.DisplayError = ""
	case ProtocolsLoadFailed:
		pvStore.Loading = false
		pvStore.DisplayError = "Protocols: " + string(e)
	default:
		return
	}
	core.Dispatch(ChangeEvent{pvStore})
}

// actions

type ProtocolsUseSocket bool

type ProtocolsLoad struct {
	Force     bool
	UseSocket bool
}

type ProtocolsLoaded *api.ProtocolInfoResult

type ProtocolsLoadFailed string

func protocolsLoad(force, useSocket bool) bool {
	return pvStore.Protocols.Query(
		useSocket, force,
		&api.Query{Type: api.QueryProtocolInfo},
		func(r *api.Query) (*api.ProtocolInfoResult, string) {
			if r.Payload == nil {
				return &api.ProtocolInfoResult{}, ""
			}
			if p, ok := r.Payload.(*api.ProtocolInfoResult); ok {
				return p, ""
			}
			return nil, "Unexpected response type"
		},
		func(v *api.ProtocolInfoResult) {
			core.Dispatch(ProtocolsLoaded(v))
		},
		func(v string) {
			core.Dispatch(ProtocolsLoadFailed(v))
		},
	)
}
