package actions

import (
	"strings"

	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/page/core"
)

func init() {
	core.DispatcherSubscribe(dvAction)
}

// store

type DiscoveryViewStore struct {
	Loading      bool
	Protocols    *api.ProtocolInfoResult
	Discoveries  map[uuid.UUID]*DiscoveryData
	DisplayError string
	sendTimeout  *core.Timeout
}

var dvStore = &DiscoveryViewStore{
	Loading:     true,
	Discoveries: make(map[uuid.UUID]*DiscoveryData),
	sendTimeout: &core.Timeout{},
}

func GetDiscoveryViewStore() *DiscoveryViewStore {
	return dvStore
}

// reducer

func dvAction(event interface{}) {
	switch e := event.(type) {
	case DiscoveryLoad:
		loadDiscoveries()
		protocolsLoadWithMainSocket()
	case ProtocolsLoaded:
		dvStore.Protocols = e
		queryDiscoveryStatus()
	case ProtocolsLoadFailed:
		pvStore.Loading = true
		// TODO
	case core.MainSocketOpened:
		pvStore.Loading = false
	case core.MainSocketMessage:
		switch e.Type {
		case api.QueryProtocolDiscoveryResult:
			if _, ok := e.Payload.(*api.ProtocolDiscoveryResult); ok {
				// TODO
			}
			queryDiscoveryStatus()
		case api.QueryProtocolDiscoveryStarted:
			// TODO
		case api.QueryProtocolDiscoveryFinished:
			// TODO
		default:
			return
		}
	case core.MainSocketError:
		pvStore.Loading = true
		// TODO
	}
	core.Dispatch(ChangeEvent{dvStore})
}

// actions

type DiscoveryLoad struct{}

type DiscoveryData struct {
	Input  *api.ProtocolDiscover
	Result *api.ProtocolDiscoveryResult
}

const localStorageDiscoveryKey = "hw-discoveries"

func loadDiscoveries() {
	if value, ok := core.LocalStorageGet(localStorageDiscoveryKey); ok {
		for _, idString := range strings.Fields(value) {
			if id, err := uuid.Parse(idString); err == nil {
				if _, ok := dvStore.Discoveries[id]; !ok {
					dvStore.Discoveries[id] = &DiscoveryData{}
				}
			}
		}
	}
}

func saveDiscoveries() {
	var sb strings.Builder
	for id := range dvStore.Discoveries {
		if sb.Len() > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(id.String())
	}
	core.LocalStorageSet(localStorageDiscoveryKey, sb.String())
}

func queryDiscoveryStatus() {
	for id, data := range dvStore.Discoveries {
		if data.Result == nil {
			core.MainSocket().Send(&api.Query{Type: api.QueryProtocolDiscovery, Payload: &api.ProtocolDiscovery{ID: id}})
			break
		}
	}
}
