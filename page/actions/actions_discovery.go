package actions

import (
	"strings"

	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/page/core"
)

func init() {
	core.DispatcherSubscribe(GetDiscoveryViewStore().action)
}

// store

type DiscoveryViewStore struct {
	Loading      bool
	Initialized  bool
	Protocols    *api.ProtocolInfoResult
	Discoveries  map[uuid.UUID]*DiscoveryData
	DisplayError string
	initTimeout  *core.Timeout
	opTimeout    *core.Timeout
}

var dvStore = &DiscoveryViewStore{
	Loading:     true,
	Initialized: false,
	Discoveries: make(map[uuid.UUID]*DiscoveryData),
	initTimeout: &core.Timeout{},
	opTimeout:   &core.Timeout{},
}

func GetDiscoveryViewStore() *DiscoveryViewStore {
	return dvStore
}

// reducer

func (s *DiscoveryViewStore) action(event interface{}) {
	switch e := event.(type) {
	case DiscoveryLoad:
		if s.Initialized {
			return
		}
		s.Initialized = true
		s.loadDiscoveries()
		GetProtocolViewStore().protocolsLoadWithMainSocket()
	case ProtocolsLoaded:
		s.Protocols = e
		s.queryDiscoveryStatus()
	case ProtocolsLoadFailed:
		s.Loading = true
		s.Initialized = false
		// TODO
	case core.MainSocketOpened:
		s.Loading = false
	case core.MainSocketMessage:
		switch e.Type {
		case api.QueryProtocolDiscoveryResult:
			s.initTimeout.Clear()
			if p, ok := e.Payload.(*api.ProtocolDiscoveryResult); ok {
				if p.Error != nil && p.Error.Code == api.ErrorNoDiscoveryID {
					delete(s.Discoveries, p.ID)
				} else {
					s.findOrAddDiscovery(p.ID).Result = p
				}
			}
			s.queryDiscoveryStatus()
		case api.QueryProtocolDiscoveryStarted:
			if p, ok := e.Payload.(*api.ProtocolDiscoveryStarted); ok {
				s.findOrAddDiscovery(p.ID).Input = &p.ProtocolDiscover
			}
		case api.QueryProtocolDiscoveryFinished:
			if p, ok := e.Payload.(*api.ProtocolDiscoveryResult); ok {
				s.findOrAddDiscovery(p.ID).Result = p
			}
		default:
			return
		}
	case core.MainSocketError:
		s.Loading = true
		s.Initialized = false
		s.initTimeout.Clear()
		s.opTimeout.Clear()
		// TODO
	case core.MainSocketTimeout:
		switch uint(e) {
		case s.initTimeout.ID:
			s.initTimeout.Clear()
			// TODO
		case s.opTimeout.ID:
			s.opTimeout.Clear()
			// TODO
		}
	}
	core.Dispatch(ChangeEvent{s})
}

// actions

type DiscoveryLoad struct{}

type DiscoveryData struct {
	Input  *api.ProtocolDiscover
	Result *api.ProtocolDiscoveryResult
}

const localStorageDiscoveryKey = "hw-discoveries"

func (s *DiscoveryViewStore) findOrAddDiscovery(id uuid.UUID) (r *DiscoveryData) {
	if r, ok := s.Discoveries[id]; !ok {
		r = &DiscoveryData{}
		s.Discoveries[id] = r
	}
	return
}

func (s *DiscoveryViewStore) loadDiscoveries() {
	if value, ok := core.LocalStorageGet(localStorageDiscoveryKey); ok {
		for _, idString := range strings.Fields(value) {
			if id, err := uuid.Parse(idString); err == nil {
				s.findOrAddDiscovery(id)
			}
		}
	}
}

func (s *DiscoveryViewStore) saveDiscoveries() {
	var sb strings.Builder
	for id := range s.Discoveries {
		if sb.Len() > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(id.String())
	}
	core.LocalStorageSet(localStorageDiscoveryKey, sb.String())
}

func (s *DiscoveryViewStore) queryDiscoveryStatus() {
	for id, data := range s.Discoveries {
		if data.Result == nil {
			core.MainSocket().SendWithTimeout(
				&api.Query{Type: api.QueryProtocolDiscovery, Payload: &api.ProtocolDiscovery{ID: id}},
				s.initTimeout,
			)
			break
		}
	}
}
