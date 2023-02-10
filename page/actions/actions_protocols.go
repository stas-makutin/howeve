package actions

import (
	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/page/core"
)

func init() {
	core.DispatcherSubscribe(GetProtocolViewStore().action)
}

// store

type ProtocolViewStore struct {
	Loading       bool
	UseSocket     bool
	Protocols     core.CachedQuery[api.ProtocolInfoResult]
	DisplayError  string
	SocketTimeout *core.Timeout
}

var pvStore = &ProtocolViewStore{
	Loading:       true,
	UseSocket:     true,
	SocketTimeout: &core.Timeout{},
}

func GetProtocolViewStore() *ProtocolViewStore {
	return pvStore
}

// reducer

func (s *ProtocolViewStore) action(event interface{}) {
	switch e := event.(type) {
	case ProtocolsUseSocket:
		s.UseSocket = bool(e)
	case *ProtocolsLoad:
		s.Loading = true
		s.DisplayError = ""
		if s.protocolsLoad(e.Force, e.UseSocket) {
			return
		}
	case ProtocolsLoaded:
		s.Loading = false
		s.DisplayError = ""
	case ProtocolsLoadFailed:
		s.Loading = false
		s.DisplayError = "Protocols: " + string(e)
	case core.MainSocketMessage:
		if e.Type == api.QueryProtocolInfoResult {
			if p, ok := e.Payload.(*api.ProtocolInfoResult); ok {
				core.Dispatch(ProtocolsLoaded(p))
			}
		}
		return
	default:
		return
	}
	core.Dispatch(ChangeEvent{s})
}

// actions

type ProtocolsUseSocket bool

type ProtocolsLoad struct {
	Force     bool
	UseSocket bool
}

type ProtocolsLoaded *api.ProtocolInfoResult

type ProtocolsLoadFailed string

func (s *ProtocolViewStore) protocolsLoad(force, useSocket bool) bool {
	return s.Protocols.Query(
		useSocket, force,
		&api.Query{Type: api.QueryProtocolInfo},
		func(r *api.Query) (*api.ProtocolInfoResult, string) {
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

func (s *ProtocolViewStore) protocolsLoadWithMainSocket() (string, bool) {
	if s.Protocols.Value != nil {
		core.Dispatch(ProtocolsLoaded(s.Protocols.Value))
		return "", true
	}
	return core.MainSocket().SendWithTimeout(&api.Query{Type: api.QueryProtocolInfo}, s.SocketTimeout)
}
