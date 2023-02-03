package actions

import (
	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/page/core"
)

func init() {
	core.DispatcherSubscribe(dvAction)
}

// store

type DiscoveryViewStore struct {
	Protocols    *api.ProtocolInfoResult
	DisplayError string
	sendTimeout  *core.Timeout
}

var dvStore = &DiscoveryViewStore{sendTimeout: &core.Timeout{}}

func GetDiscoveryViewStore() *DiscoveryViewStore {
	return dvStore
}

// reducer

func dvAction(event interface{}) {
	switch e := event.(type) {
	case DiscoveryLoad:

	case ProtocolsLoaded:
		dvStore.Protocols = e
	case ProtocolsLoadFailed:
	}
}

// actions

type DiscoveryLoad struct{}
