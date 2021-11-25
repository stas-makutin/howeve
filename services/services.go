package services

import (
	"errors"
	"sync"

	"github.com/stas-makutin/howeve/defs"
)

var services *servicesRegistry

// errors

// ErrServiceExists is the error in case if service already exists
var ErrServiceExists error = errors.New("the service already exists")

// ErrAliasExists is the error in case if service already exists
var ErrAliasExists error = errors.New("the service alias already exists")

// public interface

// AddService adds new service
func AddService(entry *defs.ServiceEntry, alias string) error {
	return services.Add(entry, alias)
}

type serviceInfo struct {
	service defs.Service
	alias   string
	params  defs.ParamValues
}

// servicesRegistry - registry of available services
type servicesRegistry struct {
	sync.Mutex
	services map[defs.ServiceKey]*serviceInfo
	aliases  map[string]*serviceInfo
}

func newServicesRegistry() *servicesRegistry {
	return &servicesRegistry{
		services: make(map[defs.ServiceKey]*serviceInfo),
		aliases:  make(map[string]*serviceInfo),
	}
}

func (sr *servicesRegistry) Stop() {
	sr.Lock()
	defer sr.Unlock()

	for _, si := range sr.services {
		si.service.Stop()
	}
	sr.services = nil
	sr.aliases = nil
}

func (sr *servicesRegistry) Add(entry *defs.ServiceEntry, alias string) error {
	sr.Lock()
	defer sr.Unlock()

	if _, ok := sr.services[entry.Key]; ok {
		return ErrServiceExists
	}
	if _, ok := sr.aliases[alias]; ok {
		return ErrAliasExists
	}

	serviceFunc := Protocols[entry.Key.Protocol].Transports[entry.Key.Transport].ServiceFunc
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
