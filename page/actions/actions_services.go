package actions

import (
	"fmt"
	"strings"

	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/page/core"
)

func init() {
	core.DispatcherSubscribe(svAction)
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

func svAction(event interface{}) {
	var errorBuilder strings.Builder

	decreaseLoadingCount := func() {
		if svStore.Loading > 0 {
			svStore.Loading -= 1
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
		svStore.UseSocket = bool(e)
	case *ServicesLoad:
		svStore.Loading = 2
		svStore.Protocols.Error = ""
		svStore.Services.Error = ""
		cachedProtocols := protocolsLoad(e.Force, e.UseSocket)
		cachedServices := servicesLoad(e.Force, e.UseSocket)
		if cachedProtocols && cachedServices {
			return
		}
	case ProtocolsLoaded:
		decreaseLoadingCount()
		svStore.Protocols.Value = e
		svStore.Protocols.Error = ""
	case ProtocolsLoadFailed:
		decreaseLoadingCount()
		svStore.Protocols.Value = nil
		svStore.Protocols.Error = string(e)
	case ServicesLoaded:
		decreaseLoadingCount()
		core.ArrangeServices(svStore.Services.Value.Services)
	case ServicesLoadFailed:
		decreaseLoadingCount()
	case *ServicesAdd:
		svStore.Loading = 1
		serviceOpRun(svStore.UseSocket,
			&api.Query{
				Type: api.QueryAddService,
				Payload: &api.ServiceEntry{
					ServiceKey: &e.Service.ServiceKey, Params: e.Service.Params.RawParams(), Alias: e.Service.Alias,
				},
			},
		)
	case *ServicesChangeAlias:
		svStore.Loading = 1
		serviceOpRun(svStore.UseSocket,
			&api.Query{
				Type: api.QueryChangeServiceAlias,
				Payload: &api.ChangeServiceAlias{
					ServiceID: &api.ServiceID{ServiceKey: e.Service}, NewAlias: e.NewAlias,
				},
			},
		)
	case *ServicesRemove:
		svStore.Loading = 1
		serviceOpRun(svStore.UseSocket,
			&api.Query{
				Type:    api.QueryRemoveService,
				Payload: &api.ServiceID{ServiceKey: e.Service},
			},
		)
	case ServicesOpRetry:
		svStore.Loading = 1
		if serviceOpRun(svStore.UseSocket, svStore.LastOp) {
			core.Dispatch(&ServicesLoad{Force: true, UseSocket: svStore.UseSocket})
		}
	case ServicesOpSucceeded:
		decreaseLoadingCount()
		svStore.LastOp = nil
		svStore.OpError = ""
		core.Dispatch(&ServicesLoad{Force: true, UseSocket: svStore.UseSocket})
	case ServicesOpFailed:
		decreaseLoadingCount()
		svStore.OpError = string(e)
		core.Dispatch(&ServicesLoad{Force: true, UseSocket: svStore.UseSocket})
	default:
		return
	}

	appendDisplayError(svStore.OpError, "")
	appendDisplayError(svStore.Protocols.Error, "Protocols: ")
	appendDisplayError(svStore.Services.Error, "Services: ")
	svStore.DisplayError = errorBuilder.String()

	core.Dispatch(ChangeEvent{svStore})
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

func servicesLoad(force, useSocket bool) bool {
	return svStore.Services.Query(
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

func serviceOpRun(useSocket bool, query *api.Query) bool {
	if query != nil {
		svStore.LastOp = query
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
