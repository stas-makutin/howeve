package services

import (
	"sync"

	"github.com/stas-makutin/howeve/defs"
)

// Service interface, defines minimal set of methods the service needs to support
type Service interface {
	Start() error
	Stop()
	Send(message defs.Message) error
}

// ServicesRegistry - registry of available services
type ServicesRegistry struct {
	sync.Mutex
	services map[defs.ServiceKey]*Service
	aliases  map[string]*defs.ServiceKey
}

func newServicesRegistry() *ServicesRegistry {
	return &ServicesRegistry{
		services: make(map[defs.ServiceKey]*Service),
		aliases:  make(map[string]*defs.ServiceKey),
	}
}

func (sr *ServicesRegistry) Stop() {
	sr.Lock()
	defer sr.Unlock()

	for _, svc := range sr.services {
		(*svc).Stop()
	}
	sr.services = nil
	sr.aliases = nil
}

func (sr *ServicesRegistry) Add(entry *defs.ServiceEntry) {
}
