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
	Initialized  bool
	Protocols    *api.ProtocolInfoResult
	Discoveries  map[uuid.UUID]*DiscoveryData
	DisplayError string
	sendTimeout  *core.Timeout
}

var dvStore = &DiscoveryViewStore{
	Loading:     true,
	Initialized: false,
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
		if dvStore.Initialized {
			return
		}
		dvStore.Initialized = true
		loadDiscoveries()
		protocolsLoadWithMainSocket()
	case ProtocolsLoaded:
		dvStore.Protocols = e
		queryDiscoveryStatus()
	case ProtocolsLoadFailed:
		dvStore.Loading = true
		dvStore.Initialized = false
		// TODO
	case core.MainSocketOpened:
		dvStore.Loading = false
	case core.MainSocketMessage:
		switch e.Type {
		case api.QueryProtocolDiscoveryResult:
			if p, ok := e.Payload.(*api.ProtocolDiscoveryResult); ok {
				if p.Error != nil && p.Error.Code == api.ErrorNoDiscoveryID {
					delete(dvStore.Discoveries, p.ID)
				} else {
					findOrAddDiscovery(p.ID).Result = p
				}
			}
			queryDiscoveryStatus()
		case api.QueryProtocolDiscoveryStarted:
			if p, ok := e.Payload.(*api.ProtocolDiscoveryStarted); ok {
				findOrAddDiscovery(p.ID).Input = &p.ProtocolDiscover
			}
		case api.QueryProtocolDiscoveryFinished:
			if p, ok := e.Payload.(*api.ProtocolDiscoveryResult); ok {
				findOrAddDiscovery(p.ID).Result = p
			}
		default:
			return
		}
	case core.MainSocketError:
		dvStore.Loading = true
		dvStore.Initialized = false
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

func findOrAddDiscovery(id uuid.UUID) (r *DiscoveryData) {
	if r, ok := dvStore.Discoveries[id]; !ok {
		r = &DiscoveryData{}
		dvStore.Discoveries[id] = r
	}
	return
}

func loadDiscoveries() {
	if value, ok := core.LocalStorageGet(localStorageDiscoveryKey); ok {
		for _, idString := range strings.Fields(value) {
			if id, err := uuid.Parse(idString); err == nil {
				findOrAddDiscovery(id)
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
