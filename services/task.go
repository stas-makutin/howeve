package services

import (
	"errors"
	"sync"

	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/defs"
	"github.com/stas-makutin/howeve/log"
	"github.com/stas-makutin/howeve/tasks"
)

// log constants
const (
	// operation
	svcOpAddFromConfig = "C"

	// operation codes
	svcOcSuccess                  = "0"
	svcOcCfgUnknownProtocol       = "P"
	svcOcCfgUnknownTransport      = "T"
	svcOcCfgNoProtocolTransport   = "X"
	svcOcCfgUnknownParameter      = "N"
	svcOcCfgNoRequiredParameter   = "R"
	svcOcCfgInvalidParameterValue = "V"
	svcOcCfgAliasExists           = "A"
	svcOcCfgCreateError           = "C"
)

// discovery constants
const (
	discoveryMaxCount  = 10
	discoveryMaxActive = 3
)

type serviceInfo struct {
	service defs.Service
	key     *defs.ServiceKey
	alias   string
	params  defs.ParamValues
}

// servicesRegistry - registry of available services - services Task implementation
type servicesRegistry struct {
	lock     sync.Mutex
	services map[defs.ServiceKey]*serviceInfo
	aliases  map[string]*serviceInfo

	cfg []config.ServiceConfig

	*discoveryRegistry
}

// NewTask func
func NewTask() *servicesRegistry {
	sr := &servicesRegistry{}
	config.AddReader(sr.readConfig)
	config.AddWriter(sr.writeConfig)
	return sr
}

func (sr *servicesRegistry) readConfig(cfg *config.Config, cfgError config.Error) {
	sr.cfg = cfg.Services
}

func (sr *servicesRegistry) writeConfig(cfg *config.Config) {
	sr.lock.Lock()
	defer sr.lock.Unlock()
	cfg.Services = sr.cfg
}

func (sr *servicesRegistry) Open(ctx *tasks.ServiceTaskContext) error {
	defs.Services = sr
	sr.services = make(map[defs.ServiceKey]*serviceInfo)
	sr.aliases = make(map[string]*serviceInfo)

	sr.discoveryRegistry = newDiscoveryRegistry(discoveryMaxCount, discoveryMaxActive)

	sr.addFromConfig()

	return nil
}

func (sr *servicesRegistry) Close(ctx *tasks.ServiceTaskContext) error {
	defs.Services = nil
	sr.services = nil
	sr.aliases = nil
	sr.discoveryRegistry = nil
	return nil
}

func (sr *servicesRegistry) Stop(ctx *tasks.ServiceTaskContext) {
	sr.lock.Lock()
	defer sr.lock.Unlock()
	defer sr.discoveryRegistry.stop()

	for _, si := range sr.services {
		si.service.Stop()
	}
}

func (sr *servicesRegistry) Add(key *defs.ServiceKey, params defs.RawParamValues, alias string) error {
	updateConfig := false
	defer func() {
		if updateConfig {
			config.WriteConfig(false)
		}
	}()
	sr.lock.Lock()
	defer sr.lock.Unlock()

	if err := sr.add(key, params, alias); err != nil {
		return err
	}

	sr.cfg = append(sr.cfg, config.ServiceConfig{
		Alias:     alias,
		Protocol:  defs.ProtocolName(key.Protocol),
		Transport: defs.TransportName(key.Transport),
		Entry:     key.Entry,
		Params:    params,
	})
	updateConfig = true

	return nil
}

// Alias changes service's alias
func (sr *servicesRegistry) Alias(key *defs.ServiceKey, alias string) error {
	updateConfig := false
	defer func() {
		if updateConfig {
			config.WriteConfig(false)
		}
	}()
	sr.lock.Lock()
	defer sr.lock.Unlock()

	var si *serviceInfo
	if key != nil {
		si = sr.services[*key]
	}
	if si == nil {
		return defs.ErrServiceNotExists
	}

	if si.alias != alias {
		if si.alias != "" {
			delete(sr.aliases, si.alias)
		}
		si.alias = alias
		if alias != "" {
			sr.aliases[alias] = si
		}

		if i, ok := sr.findServiceCfg(si.key); ok {
			sr.cfg[i].Alias = si.alias
			updateConfig = true
		}
	}

	return nil
}

// Remove removes the service identified by (in order of priority): 1) service key; 2) alias
func (sr *servicesRegistry) Remove(key *defs.ServiceKey, alias string) error {
	updateConfig := false
	defer func() {
		if updateConfig {
			config.WriteConfig(false)
		}
	}()
	sr.lock.Lock()
	defer sr.lock.Unlock()

	var si *serviceInfo
	if key != nil {
		si = sr.services[*key]
	} else if alias != "" {
		si = sr.aliases[alias]
	}
	if si == nil {
		return defs.ErrServiceNotExists
	}

	si.service.Stop()

	delete(sr.services, *si.key)
	if si.alias != "" {
		delete(sr.aliases, si.alias)
	}

	if i, ok := sr.findServiceCfg(si.key); ok {
		sr.cfg = append(sr.cfg[:i], sr.cfg[i+1:]...)
		updateConfig = true
	}

	return nil
}

// Status return the status of the service identified by (in order of priority): 1) service key; 2) alias
func (sr *servicesRegistry) Status(key *defs.ServiceKey, alias string) (defs.ServiceStatus, bool) {
	sr.lock.Lock()
	defer sr.lock.Unlock()

	var si *serviceInfo
	if key != nil {
		si = sr.services[*key]
	} else if alias != "" {
		si = sr.aliases[alias]
	}
	if si == nil {
		return nil, false
	}

	return si.service.Status(), true
}

// List returns list of registered services to provided callback function. The services iteration will stop if callback function will return true
func (sr *servicesRegistry) List(listFn defs.ListFunc) {
	if listFn != nil {
		sr.lock.Lock()
		defer sr.lock.Unlock()

		for _, si := range sr.services {
			listFn(si.key, si.alias, si.service.Status())
		}
	}
}

// Send sends payload to the service identified by (in order of priority): 1) service key; 2) alias
func (sr *servicesRegistry) Send(key *defs.ServiceKey, alias string, payload []byte) (*defs.Message, error) {
	sr.lock.Lock()
	defer sr.lock.Unlock()

	var si *serviceInfo
	if key != nil {
		si = sr.services[*key]
	} else if alias != "" {
		si = sr.aliases[alias]
	}
	if si == nil {
		return nil, defs.ErrServiceNotExists
	}

	return si.service.Send(payload)
}

func (sr *servicesRegistry) add(key *defs.ServiceKey, params defs.RawParamValues, alias string) error {
	if _, ok := sr.services[*key]; ok {
		return defs.ErrServiceExists
	}
	if _, ok := sr.aliases[alias]; ok {
		return defs.ErrAliasExists
	}

	to, ti, err := defs.ResolveProtocolAndTransport(key.Protocol, key.Transport)
	if err != nil {
		return err
	}

	serviceFunc := to.ServiceFunc
	if serviceFunc == nil {
		return defs.ErrNoProtocolTransport
	}

	pv, err := ti.Params.Merge(to.Params).ParseValues(params)
	if err != nil {
		return err
	}

	service, error := serviceFunc(key.Entry, pv)
	if error != nil {
		return error
	}

	service.Start()

	si := &serviceInfo{service, key, alias, pv}
	sr.services[*key] = si
	if alias != "" {
		sr.aliases[alias] = si
	}
	return nil
}

func (sr *servicesRegistry) addFromConfig() {
	for _, cfg := range sr.cfg {
		protocol, ok := defs.ProtocolByName(cfg.Protocol)
		if !ok {
			log.Report(log.SrcSVC, svcOpAddFromConfig, svcOcCfgUnknownProtocol, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias)
			continue
		}
		transport, ok := defs.TransportByName(cfg.Transport)
		if !ok {
			log.Report(log.SrcSVC, svcOpAddFromConfig, svcOcCfgUnknownTransport, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias)
			continue
		}

		err := sr.add(
			&defs.ServiceKey{
				Protocol:  protocol,
				Transport: transport,
				Entry:     cfg.Entry,
			},
			cfg.Params,
			cfg.Alias,
		)

		var pe *defs.ParseError
		switch {
		case errors.Is(err, defs.ErrServiceExists):
			// ignore
		case errors.Is(err, defs.ErrAliasExists):
			log.Report(log.SrcSVC, svcOpAddFromConfig, svcOcCfgAliasExists, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias)
		case errors.Is(err, defs.ErrProtocolNotSupported):
			log.Report(log.SrcSVC, svcOpAddFromConfig, svcOcCfgUnknownProtocol, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias)
		case errors.Is(err, defs.ErrTransportNotSupported):
			log.Report(log.SrcSVC, svcOpAddFromConfig, svcOcCfgUnknownTransport, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias)
		case errors.Is(err, defs.ErrNoProtocolTransport):
			log.Report(log.SrcSVC, svcOpAddFromConfig, svcOcCfgNoProtocolTransport, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias)
		case errors.As(err, &pe):
			switch pe.Code {
			case defs.UnknownParamName:
				log.Report(log.SrcSVC, svcOpAddFromConfig, svcOcCfgUnknownParameter, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias, pe.Name)
			case defs.NoRequiredParam:
				log.Report(log.SrcSVC, svcOpAddFromConfig, svcOcCfgNoRequiredParameter, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias, pe.Name)
			default:
				log.Report(log.SrcSVC, svcOpAddFromConfig, svcOcCfgInvalidParameterValue, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias, pe.Name, pe.Value)
			}
		case err != nil:
			log.Report(log.SrcSVC, svcOpAddFromConfig, svcOcCfgCreateError, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias, err.Error())
		}
	}
}

func (sr *servicesRegistry) findServiceCfg(key *defs.ServiceKey) (int, bool) {
	protocolName := defs.ProtocolName(key.Protocol)
	transportName := defs.TransportName(key.Transport)

	for i, cfg := range sr.cfg {
		if cfg.Transport == transportName && cfg.Protocol == protocolName && cfg.Entry == key.Entry {
			return i, true
		}
	}
	return 0, false
}
