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
		core.Dispatch(ServicesLoadFailed("Unexpected response type"))
	}
}

func servicesLoadSockets() {
	core.FetchQueryWithSocket(
		&api.Query{Type: api.QueryListServices},
		func(r *api.Query) {
			servicesListProcessResponse(r)
		},
		func(err string) {
			core.Dispatch(ServicesLoadFailed(err))
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

func servicesLoad(action *ServicesLoad) bool {
	if action.Force || (svStore.Services == nil && svStore.Error == "") {
		if action.UseSocket {
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
	Loading   bool
	UseSocket bool
	Error     string
	Services  *api.ListServicesResult
}

var svStore = &ServicesViewStore{
	Loading:   true,
	UseSocket: true,
}
var svStoreChanging = false

func GetServicesViewStore() *ServicesViewStore {
	return svStore
}

// reducer

func svAction(event interface{}) {
	switch e := event.(type) {
	case ServicesUseSocket:
		svStore.UseSocket = bool(e)
	case *ServicesLoad:
		svStore.Loading = true
		svStore.Error = ""
		if servicesLoad(e) {
			return
		}
	case *ServicesLoaded:
		svStore.Loading = false
		svStore.Error = ""
		svStore.Services = e.Services
	case ServicesLoadFailed:
		svStore.Loading = false
		svStore.Error = string(e)
		if svStore.Error == "" {
			svStore.Error = "Unable to load Services"
		}
	default:
		return
	}
	core.Dispatch(ChangeEvent{svStore})
}
