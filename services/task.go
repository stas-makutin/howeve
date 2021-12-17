package services

import (
	"errors"
	"sync"

	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/defs"
	"github.com/stas-makutin/howeve/log"
	"github.com/stas-makutin/howeve/tasks"
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

	sr.addFromConfig()

	return nil
}

func (sr *servicesRegistry) Close(ctx *tasks.ServiceTaskContext) error {
	defs.Services = nil
	sr.services = nil
	sr.aliases = nil
	return nil
}

func (sr *servicesRegistry) Stop(ctx *tasks.ServiceTaskContext) {
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

	pi := defs.Protocols[entry.Key.Protocol]
	if pi == nil {
		return defs.ErrProtocolNotSupported
	}
	to := pi.Transports[entry.Key.Transport]
	if to == nil {
		return defs.ErrTransportNotSupported
	}
	serviceFunc := to.ServiceFunc
	if serviceFunc == nil {
		return defs.ErrTransportNotSupported
	}

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

func (sr *servicesRegistry) addFromConfig() {
	for _, cfg := range sr.cfg {
		protocol, ok := defs.ProtocolByName(cfg.Protocol)
		if !ok {
			log.Report(log.SrcSVC, SvcAddFromConfig, SvcOcCfgUnknownProtocol, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias)
			continue
		}
		transport, ok := defs.TransportByName(cfg.Transport)
		if !ok {
			log.Report(log.SrcSVC, SvcAddFromConfig, SvcOcCfgUnknownTransport, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias)
			continue
		}
		pi, ok := defs.Protocols[protocol]
		if !ok {
			log.Report(log.SrcSVC, SvcAddFromConfig, SvcOcCfgProtocolNotSupported, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias)
			continue
		}
		pto, ok := pi.Transports[transport]
		if !ok {
			log.Report(log.SrcSVC, SvcAddFromConfig, SvcOcCfgTransportNotSupported, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias)
			continue
		}
		ti, ok := defs.Transports[transport]
		if !ok {
			log.Report(log.SrcSVC, SvcAddFromConfig, SvcOcCfgTransportNotSupported, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias)
			continue
		}
		params, paramName, err := pto.Params.Merge(ti.Params).ParseAll(cfg.Params)
		if err != nil {
			switch {
			case errors.Is(err, defs.ErrUnknownParamName):
				log.Report(log.SrcSVC, SvcAddFromConfig, SvcOcCfgUnknownParameter, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias, paramName)
				continue
			case errors.Is(err, defs.ErrNoRequiredParam):
				log.Report(log.SrcSVC, SvcAddFromConfig, SvcOcCfgNoRequiredParameter, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias, paramName)
				continue
			}

			value := cfg.Params[paramName]
			log.Report(log.SrcSVC, SvcAddFromConfig, SvcOcCfgInvalidParameterValue, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias, paramName, value)
			continue
		}

		err = sr.Add(
			&defs.ServiceEntry{
				Key: defs.ServiceKey{
					Protocol:  protocol,
					Transport: transport,
					Entry:     cfg.Entry,
				},
				Params: params,
			},
			cfg.Alias,
		)
		switch {
		case errors.Is(err, defs.ErrServiceExists):
			// ignore
		case errors.Is(err, defs.ErrAliasExists):
			log.Report(log.SrcSVC, SvcAddFromConfig, SvcOcCfgAliasExists, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias)
		case errors.Is(err, defs.ErrProtocolNotSupported):
			log.Report(log.SrcSVC, SvcAddFromConfig, SvcOcCfgProtocolNotSupported, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias)
		case errors.Is(err, defs.ErrTransportNotSupported):
			log.Report(log.SrcSVC, SvcAddFromConfig, SvcOcCfgTransportNotSupported, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias)
		case err != nil:
			log.Report(log.SrcSVC, SvcAddFromConfig, SvcOcCfgCreateError, cfg.Protocol, cfg.Transport, cfg.Entry, cfg.Alias, err.Error())
		}
	}
}
