package actions

import (
	"encoding/json"
	"fmt"

	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/page/core"
)

func init() {
	core.DispatcherSubscribe(cvAction)
}

// actions

type ConfigUseSocket bool

type ConfigLoad struct {
	Force     bool
	UseSocket bool
}

type ConfigLoaded struct {
	Config string
}

type ConfigLoadFailed string

func configProcessResponse(r *api.Query) {
	if p, ok := r.Payload.(*api.Config); ok {
		b, err := json.MarshalIndent(p, "", "  ")
		if err == nil {
			core.Dispatch(&ConfigLoaded{string(b)})
			return
		}
	}
	core.Dispatch(ConfigLoadFailed("Unexpected response type"))
}

func configLoadSockets() {
	core.FetchQueryWithSocket(
		&api.Query{Type: api.QueryGetConfig},
		func(r *api.Query) {
			configProcessResponse(r)
		},
		func(err string) {
			core.Dispatch(ConfigLoadFailed(err))
		},
	)
}

func configLoadFetch() {
	core.FetchQuery(
		core.HTTPUrl("/cfg"), nil,
		func(r *api.Query) {
			configProcessResponse(r)
		},
		func(err string) {
			core.Dispatch(ConfigLoadFailed(err))
		},
	)
}

func configLoad(action *ConfigLoad) bool {
	if action.Force || (cvStore.Config == "" && cvStore.Error == "") {
		if action.UseSocket {
			configLoadSockets()
		} else {
			configLoadFetch()
		}
		return false
	}
	// restore saved state
	if cvStore.Error != "" {
		core.Dispatch(ConfigLoadFailed(cvStore.Error))
	} else {
		core.Dispatch(&ConfigLoaded{cvStore.Config})
	}
	return true
}

// store

type ConfigViewStore struct {
	Loading   bool
	UseSocket bool
	Error     string
	Config    string
}

var cvStore = &ConfigViewStore{
	Loading:   true,
	UseSocket: true,
}
var cvStoreChanging = false

func GetConfigViewStore() *ConfigViewStore {
	return cvStore
}

// reducer

func cvAction(event interface{}) {
	core.Console.Log(fmt.Sprintf("%T", event))
	switch e := event.(type) {
	case ConfigUseSocket:
		cvStore.UseSocket = bool(e)
	case *ConfigLoad:
		cvStore.Loading = true
		cvStore.Error = ""
		if configLoad(e) {
			return
		}
	case *ConfigLoaded:
		cvStore.Loading = false
		cvStore.Error = ""
		cvStore.Config = e.Config
	case ConfigLoadFailed:
		cvStore.Loading = false
		cvStore.Error = string(e)
		if cvStore.Error == "" {
			cvStore.Error = "Unable to load configuration"
		}
	default:
		return
	}
	core.Dispatch(ChangeEvent{cvStore})
}
