package actions

import "github.com/stas-makutin/howeve/page/core"

func init() {
	core.DispatcherSubscribe(dvAction)
}

// store

type DiscoveryViewStore struct {
}

var dvStore = &DiscoveryViewStore{}

func GetDiscoveryViewStore() *DiscoveryViewStore {
	return dvStore
}

// reducer

func dvAction(event interface{}) {
}
