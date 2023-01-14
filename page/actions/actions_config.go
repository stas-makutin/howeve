package actions

import (
	"encoding/json"

	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/page/core"
)

func init() {
	core.DispatcherSubscribe(cvAction)
}

type ConfigString struct {
	Value string
}

// store

type ConfigViewStore struct {
	Loading      bool
	UseSocket    bool
	Config       core.CachedQuery[ConfigString]
	DisplayError string
}

var cvStore = &ConfigViewStore{
	Loading:   true,
	UseSocket: true,
}

func GetConfigViewStore() *ConfigViewStore {
	return cvStore
}

func (vs *ConfigViewStore) ConfigValue() (v string) {
	if vs.Config.Value != nil {
		v = vs.Config.Value.Value
	}
	return
}

// reducer

func cvAction(event interface{}) {
	switch e := event.(type) {
	case ConfigUseSocket:
		cvStore.UseSocket = bool(e)
	case *ConfigLoad:
		cvStore.Loading = true
		cvStore.DisplayError = ""
		if configLoad(e.Force, e.UseSocket) {
			return
		}
	case ConfigLoaded:
		cvStore.Loading = false
		cvStore.DisplayError = ""
	case ConfigLoadFailed:
		cvStore.Loading = false
		cvStore.DisplayError = cvStore.Config.Error
	default:
		return
	}
	core.Dispatch(ChangeEvent{cvStore})
}

// actions

type ConfigUseSocket bool

type ConfigLoad struct {
	Force     bool
	UseSocket bool
}

type ConfigLoaded struct{}

type ConfigLoadFailed struct{}

func configLoad(force, useSocket bool) bool {
	return cvStore.Config.Query(
		useSocket, force,
		&api.Query{Type: api.QueryGetConfig},
		func(r *api.Query) (*ConfigString, string) {
			if p, ok := r.Payload.(*api.Config); ok {
				b, err := json.MarshalIndent(p, "", "  ")
				if err == nil {
					return &ConfigString{Value: string(b)}, ""
				}
			}
			return nil, "Unexpected response type"
		},
		func(s string) string {
			return "Config: " + s
		},
		func(*ConfigString) {
			core.Dispatch(ConfigLoaded{})
		},
		func(string) {
			core.Dispatch(ConfigLoadFailed{})
		},
	)
}
