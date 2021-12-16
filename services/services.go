package services

import (
	"sync"

	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/defs"
)

var services *servicesRegistry

type serviceInfo struct {
	service defs.Service
	alias   string
	params  defs.ParamValues
}

// servicesRegistry - registry of available services
type servicesRegistry struct {
	lock     sync.Mutex
	services map[defs.ServiceKey]*serviceInfo
	aliases  map[string]*serviceInfo

	cfg []config.ServiceConfig
}

func newServicesRegistry() *servicesRegistry {
	return &servicesRegistry{}
}

func (sr *servicesRegistry) readConfig(cfg *config.Config, cfgError config.Error) {
	sr.cfg = cfg.Services
}

func (sr *servicesRegistry) writeConfig(cfg *config.Config) {
	cfg.Services = sr.cfg
}

func (sr *servicesRegistry) open() {
	sr.services = make(map[defs.ServiceKey]*serviceInfo)
	sr.aliases = make(map[string]*serviceInfo)
}

func (sr *servicesRegistry) close() {
	sr.services = nil
	sr.aliases = nil
}

func (sr *servicesRegistry) stop() {
	sr.lock.Lock()
	defer sr.lock.Unlock()

	for _, si := range sr.services {
		si.service.Stop()
	}
}

func (sr *servicesRegistry) Add(entry *defs.ServiceEntry, alias string) error {
	sr.lock.Lock()
	defer sr.lock.Unlock()

	if _, ok := sr.services[entry.Key]; ok {
		return defs.ErrServiceExists
	}
	if _, ok := sr.aliases[alias]; ok {
		return defs.ErrAliasExists
	}

	serviceFunc := defs.Protocols[entry.Key.Protocol].Transports[entry.Key.Transport].ServiceFunc
	service, error := serviceFunc(entry.Key.Entry, entry.Params)
	if error != nil {
		return error
	}

	service.Start()

	si := &serviceInfo{service, alias, entry.Params}
	sr.services[entry.Key] = si
	if alias != "" {
		sr.aliases[alias] = si
	}
	return nil
}
