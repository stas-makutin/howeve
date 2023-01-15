package actions

import (
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
	CRUDError    string
	DisplayError string
}

var svStore = &ServicesViewStore{
	Loading:   1,
	UseSocket: true,
}
var svStoreChanging = false

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
	appendDisplayError := func(message string) {
		if len(message) > 0 {
			if errorBuilder.Len() > 0 {
				errorBuilder.WriteString("\n")
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
	case ServicesLoaded, ServicesLoadFailed:
		decreaseLoadingCount()
	default:
		return
	}

	appendDisplayError(svStore.CRUDError)
	appendDisplayError(svStore.Protocols.Error)
	appendDisplayError(svStore.Services.Error)
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

// ServiceAdd and triggered action

type ServiceAdd struct {
	UseSocket bool
	Service   *core.ServiceEntryData
}

type ServiceAddSuccess struct{}

type ServiceAddFailed string

// ServiceRemove and triggered action

type ServiceRemove struct {
	UseSocket bool
	Service   *api.ServiceKey
}

type ServiceRemoveSuccess struct{}

type ServiceRemoveFailed string

// ServiceChangeAlias and triggered action

type ServiceChangeAlias struct {
	UseSocket bool
	Service   *api.ServiceKey
	NewAlias  string
}

type ServiceChangeSuccess struct{}

type ServiceChangeFailed string
