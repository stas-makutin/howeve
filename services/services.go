package services

import (
	"context"
	"sync"

	"github.com/stas-makutin/howeve/defs"
)

var services *servicesRegistry

// public interface

// AddService adds new service
func AddService(entry *defs.ServiceEntry, alias string) error {
	return services.Add(entry, alias)
}

type serviceInfo struct {
	service *defs.Service
	alias   string
	params  defs.ParamValues
}

// servicesRegistry - registry of available services
type servicesRegistry struct {
	sync.Mutex
	ctx      context.Context
	cancel   context.CancelFunc
	services map[defs.ServiceKey]*serviceInfo
	aliases  map[string]*defs.ServiceKey
}

func newServicesRegistry() *servicesRegistry {
	ctx, cancel := context.WithCancel(context.Background())
	return &servicesRegistry{
		ctx:      ctx,
		cancel:   cancel,
		services: make(map[defs.ServiceKey]*serviceInfo),
		aliases:  make(map[string]*defs.ServiceKey),
	}
}

func (sr *servicesRegistry) Stop() {
	sr.Lock()
	defer sr.Unlock()

	sr.cancel()
	for _, si := range sr.services {
		(*si.service).Stop()
	}
	sr.services = nil
	sr.aliases = nil
}

func (sr *servicesRegistry) Add(entry *defs.ServiceEntry, alias string) error {
	sr.Lock()
	defer sr.Unlock()

	if _, ok := sr.services[entry.Key]; ok {
		// already exists
	}
	if _, ok := sr.aliases[alias]; ok {
		// alias already exists
	}

	return nil
}
