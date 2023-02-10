package actions

import (
	"fmt"
	"strings"

	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/page/core"
)

func init() {
	core.DispatcherSubscribe(GetServicesViewStore().action)
}

// store

type ServicesViewStore struct {
	Loading      int
	UseSocket    bool
	Protocols    core.CachedQuery[api.ProtocolInfoResult]
	Services     core.CachedQuery[api.ListServicesResult]
	LastOp       *api.Query
	OpError      string
	DisplayError string
}

var svStore = &ServicesViewStore{
	Loading:   1,
	UseSocket: true,
}

func GetServicesViewStore() *ServicesViewStore {
	return svStore
}

// reducer

func (s *ServicesViewStore) action(event interface{}) {
	var errorBuilder strings.Builder

	decreaseLoadingCount := func() {
		if s.Loading > 0 {
			s.Loading -= 1
		}
	}
	appendDisplayError := func(message, prefix string) {
		if len(message) > 0 {
			if errorBuilder.Len() > 0 {
				errorBuilder.WriteString("\n")
			}
			if len(prefix) > 0 {
				errorBuilder.WriteString(prefix)
			}
			errorBuilder.WriteString(message)
		}
	}

	switch e := event.(type) {
	case ServicesUseSocket:
		s.UseSocket = bool(e)
	case *ServicesLoad:
		s.Loading = 2
		s.Protocols.Error = ""
		s.Services.Error = ""
		cachedProtocols := GetProtocolViewStore().protocolsLoad(e.Force, e.UseSocket)
		cachedServices := s.servicesLoad(e.Force, e.UseSocket)
		if cachedProtocols && cachedServices {
			return
		}
	case ProtocolsLoaded:
		decreaseLoadingCount()
		s.Protocols.Value = e
		s.Protocols.Error = ""
	case ProtocolsLoadFailed:
		decreaseLoadingCount()
		s.Protocols.Value = nil
		s.Protocols.Error = string(e)
	case ServicesLoaded:
		decreaseLoadingCount()
		core.ArrangeServices(s.Services.Value.Services)
	case ServicesLoadFailed:
		decreaseLoadingCount()
	case *ServicesAdd:
		s.Loading = 1
		s.serviceOpRun(s.UseSocket,
			&api.Query{
				Type: api.QueryAddService,
				Payload: &api.ServiceEntry{
					ServiceKey: &e.Service.ServiceKey, Params: e.Service.Params.RawParams(), Alias: e.Service.Alias,
				},
			},
		)
	case *ServicesChangeAlias:
		s.Loading = 1
		s.serviceOpRun(s.UseSocket,
			&api.Query{
				Type: api.QueryChangeServiceAlias,
				Payload: &api.ChangeServiceAlias{
					ServiceID: &api.ServiceID{ServiceKey: e.Service}, NewAlias: e.NewAlias,
				},
			},
		)
	case *ServicesRemove:
		s.Loading = 1
		s.serviceOpRun(s.UseSocket,
			&api.Query{
				Type:    api.QueryRemoveService,
				Payload: &api.ServiceID{ServiceKey: e.Service},
			},
		)
	case ServicesOpRetry:
		s.Loading = 1
		if s.serviceOpRun(s.UseSocket, s.LastOp) {
			core.Dispatch(&ServicesLoad{Force: true, UseSocket: s.UseSocket})
		}
	case ServicesOpSucceeded:
		decreaseLoadingCount()
		s.LastOp = nil
		s.OpError = ""
		core.Dispatch(&ServicesLoad{Force: true, UseSocket: s.UseSocket})
	case ServicesOpFailed:
		decreaseLoadingCount()
		s.OpError = string(e)
		core.Dispatch(&ServicesLoad{Force: true, UseSocket: s.UseSocket})
	default:
		return
	}

	appendDisplayError(s.OpError, "")
	appendDisplayError(s.Protocols.Error, "Protocols: ")
	appendDisplayError(s.Services.Error, "Services: ")
	s.DisplayError = errorBuilder.String()

	core.Dispatch(ChangeEvent{s})
}

// actions

// ServicesUseSocket action

type ServicesUseSocket bool

// ServicesLoad and triggered action

type ServicesLoad struct {
	Force     bool
	UseSocket bool
}

type ServicesLoaded *api.ListServicesResult

type ServicesLoadFailed string

func (s *ServicesViewStore) servicesLoad(force, useSocket bool) bool {
	return s.Services.Query(
		useSocket, force,
		&api.Query{Type: api.QueryListServices},
		func(r *api.Query) (*api.ListServicesResult, string) {
			if r.Payload == nil {
				return &api.ListServicesResult{}, ""
			} else if p, ok := r.Payload.(*api.ListServicesResult); ok {
				return p, ""
			} else {
				return nil, "Unexpected response type"
			}
		},
		func(v *api.ListServicesResult) {
			core.Dispatch(ServicesLoaded(v))
		},
		func(v string) {
			core.Dispatch(ServicesLoadFailed(v))
		},
	)
}

// Operations (Add Service, Change Service Alias, Remove Service) action

type ServicesOpRetry struct{}

type ServicesOpSucceeded struct{}

type ServicesOpFailed string

func (s *ServicesViewStore) serviceOpRun(useSocket bool, query *api.Query) bool {
	if query != nil {
		s.LastOp = query
		errorPrefix := ""
		switch query.Type {
		case api.QueryAddService:
			errorPrefix = "Add Service: "
		case api.QueryChangeServiceAlias:
			errorPrefix = "Change Service Alias: "
		case api.QueryRemoveService:
			errorPrefix = "Remove Service: "
		}
		core.Query(
			useSocket, query,
			func(q *api.Query) {
				if status, ok := q.Payload.(*api.StatusReply); ok {
					if status.Success {
						core.Dispatch(ServicesOpSucceeded{})
					} else {
						core.Dispatch(ServicesOpFailed(fmt.Sprintf("%s[%d] %s", errorPrefix, status.Error.Code, status.Error.Message)))
					}
					return
				}
				core.Dispatch(ServicesOpFailed(errorPrefix + "Unexpected response type"))
			},
			func(err string) {
				core.Dispatch(ServicesOpFailed(errorPrefix + err))
			},
		)
		return false
	}
	return true
}

// ServicesAdd action

type ServicesAdd struct {
	UseSocket bool
	Service   *core.ServiceEntryData
}

// ServicesRemove action

type ServicesRemove struct {
	UseSocket bool
	Service   *api.ServiceKey
}

// ServicesChangeAlias action

type ServicesChangeAlias struct {
	UseSocket bool
	Service   *api.ServiceKey
	NewAlias  string
}
