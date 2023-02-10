package actions

import (
	"encoding/json"

	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/page/core"
)

func init() {
	core.DispatcherSubscribe(GetConfigViewStore().action)
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

func (s *ConfigViewStore) action(event interface{}) {
	switch e := event.(type) {
	case ConfigUseSocket:
		s.UseSocket = bool(e)
	case *ConfigLoad:
		s.Loading = true
		s.DisplayError = ""
		if s.configLoad(e.Force, e.UseSocket) {
			return
		}
	case ConfigLoaded:
		s.Loading = false
		s.DisplayError = ""
	case ConfigLoadFailed:
		s.Loading = false
		s.DisplayError = "Config: " + string(e)
	default:
		return
	}
	core.Dispatch(ChangeEvent{s})
}

// actions

type ConfigUseSocket bool

type ConfigLoad struct {
	Force     bool
	UseSocket bool
}

type ConfigLoaded string

type ConfigLoadFailed string

func (s *ConfigViewStore) configLoad(force, useSocket bool) bool {
	return s.Config.Query(
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
		func(v *ConfigString) {
			core.Dispatch(ConfigLoaded(v.Value))
		},
		func(v string) {
			core.Dispatch(ConfigLoadFailed(v))
		},
	)
}
