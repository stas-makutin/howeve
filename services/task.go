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
	svcOcCfgProtocolNotSupported  = "X"
	svcOcCfgTransportNotSupported = "x"
	svcOcCfgUnknownParameter      = "N"
	svcOcCfgNoRequiredParameter   = "R"
	svcOcCfgInvalidParameterValue = "V"
	svcOcCfgAliasExists           = "A"
	svcOcCfgCreateError           = "C"
)

type serviceInfo struct {
	service defs.Service
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
	cfg.Services = sr.cfg
}

func (sr *servicesRegistry) Open(ctx *tasks.ServiceTaskContext) error {
	defs.Services = sr
	sr.services = make(map[defs.ServiceKey]*serviceInfo)
	sr.aliases = make(map[string]*serviceInfo)

	sr.discoveryRegistry = newDiscoveryRegistry()

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

	config.WriteConfig(false)

	return nil
}

func (sr *servicesRegistry) add(key *defs.ServiceKey, params defs.RawParamValues, alias string) error {
	if _, ok := sr.services[*key]; ok {
		return defs.ErrServiceExists
	}
	if _, ok := sr.aliases[alias]; ok {
		return defs.ErrAliasExists
	}

	pi := defs.Protocols[key.Protocol]
	if pi == nil {
		return defs.ErrProtocolNotSupported
	}
	ti := defs.Transports[key.Transport]
	if ti == nil {
		return defs.ErrTransportNotSupported
	}
	to := pi.Transports[key.Transport]
	if to == nil {
		return defs.ErrTransportNotSupported
	}
	serviceFunc := to.ServiceFunc
	if serviceFunc == nil {
		return defs.ErrTransportNotSupported
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

	si := &serviceInfo{service, alias, pv}
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
			log.Report(log.SrcSVC, svcOpAddFromConfig, svcOcCfgProtocolNotSupported, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias)
		case errors.Is(err, defs.ErrTransportNotSupported):
			log.Report(log.SrcSVC, svcOpAddFromConfig, svcOcCfgTransportNotSupported, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias)
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
