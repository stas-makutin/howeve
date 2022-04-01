package actions

import (
	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/page/core"
)

func init() {
	core.DispatcherSubscribe(svAction)
}

// actions

type ServicesUseSocket bool

type ServicesLoad struct {
	Force     bool
	UseSocket bool
}

type ServicesLoaded struct {
	Services *api.ListServicesResult
}

type ServicesLoadFailed string

func servicesListProcessResponse(r *api.Query) {
	if r.Payload == nil {
		core.Dispatch(&ServicesLoaded{&api.ListServicesResult{}})
	} else if p, ok := r.Payload.(*api.ListServicesResult); ok {
		core.Dispatch(&ServicesLoaded{p})
	} else {
		core.Dispatch(ServicesLoadFailed("Services: unexpected response type"))
	}
}

func servicesLoadSockets() {
	core.FetchQueryWithSocket(
		&api.Query{Type: api.QueryListServices},
		func(r *api.Query) {
			servicesListProcessResponse(r)
		},
		func(err string) {
			core.Dispatch(ServicesLoadFailed("Services: " + err))
		},
	)
}

func servicesLoadFetch() {
	core.FetchQuery(
		core.HTTPUrl("/cfg"), nil,
		func(r *api.Query) {
			servicesListProcessResponse(r)
		},
		func(err string) {
			core.Dispatch(ServicesLoadFailed(err))
		},
	)
}

func servicesLoad(force, useSocket bool) bool {
	if force || (svStore.Services == nil && svStore.Error == "") {
		if useSocket {
			servicesLoadSockets()
		} else {
			servicesLoadFetch()
		}
		return false
	}
	// restore saved state
	if svStore.Error != "" {
		core.Dispatch(ServicesLoadFailed(svStore.Error))
	} else {
		core.Dispatch(&ServicesLoaded{svStore.Services})
	}
	return true
}

// store

type ServicesViewStore struct {
	Loading   int
	UseSocket bool
	Error     string
	Protocols *api.ProtocolInfoResult
	Services  *api.ListServicesResult
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
	decreaseLoadingCount := func() {
		if svStore.Loading > 0 {
			svStore.Loading -= 1
		}
	}
	appendErrorMessage := func(msg, msgDef string) {
		if msg == "" {
			msg = msgDef
		}
		if msg != "" {
			if svStore.Error != "" {
				svStore.Error += "; "
			}
			svStore.Error += msg
		}
	}

	switch e := event.(type) {
	case ServicesUseSocket:
		svStore.UseSocket = bool(e)
	case *ServicesLoad:
		svStore.Loading = 2
		svStore.Error = ""
		sf := servicesLoad(e.Force, e.UseSocket)
		pf := protocolsLoad(e.Force, e.UseSocket)
		if sf && pf {
			return
		}
	case *ServicesLoaded:
		decreaseLoadingCount()
		svStore.Services = e.Services
	case *ProtocolsLoaded:
		decreaseLoadingCount()
		svStore.Protocols = e.Protocols
	case ServicesLoadFailed:
		decreaseLoadingCount()
		appendErrorMessage(string(e), "Services: unable to load")
	case ProtocolsLoadFailed:
		decreaseLoadingCount()
		appendErrorMessage(string(e), "Protocols: unable to load")
	default:
		return
	}
	core.Dispatch(ChangeEvent{svStore})
}
