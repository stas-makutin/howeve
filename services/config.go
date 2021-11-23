package services

import (
	"errors"

	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/defs"
	"github.com/stas-makutin/howeve/log"
)

func addServiceFromConfig(cfg config.ServiceConfig) bool {
	protocol, ok := ProtocolByName(cfg.Protocol)
	if !ok {
		log.Report(log.SrcSVC, SvcOpStart, SvcOcCfgUnknownProtocol, cfg.Protocol, cfg.Transport, cfg.Entry)
		return false
	}
	transport, ok := TransportByName(cfg.Transport)
	if !ok {
		log.Report(log.SrcSVC, SvcOpStart, SvcOcCfgUnknownTransport, cfg.Protocol, cfg.Transport, cfg.Entry)
		return false
	}
	pi, ok := Protocols[protocol]
	if !ok {
		log.Report(log.SrcSVC, SvcOpStart, SvcOcCfgProtocolNotSupported, cfg.Protocol, cfg.Transport, cfg.Entry)
		return false
	}
	pto, ok := pi.Transports[transport]
	if !ok {
		log.Report(log.SrcSVC, SvcOpStart, SvcOcCfgTransportNotSupported, cfg.Protocol, cfg.Transport, cfg.Entry)
		return false
	}
	ti, ok := Transports[transport]
	if !ok {
		log.Report(log.SrcSVC, SvcOpStart, SvcOcCfgTransportNotSupported, cfg.Protocol, cfg.Transport, cfg.Entry)
		return false
	}
	_ /*params*/, paramName, err := pto.Params.Merge(ti.Params).ParseAll(cfg.Params)
	if err != nil {
		if errors.Is(err, defs.ErrUnknownParamName) {
			log.Report(log.SrcSVC, SvcOpStart, SvcOcCfgUnknownParameter, cfg.Protocol, cfg.Transport, cfg.Entry, paramName)
			return false
		} else if errors.Is(err, defs.ErrNoRequiredParam) {
			log.Report(log.SrcSVC, SvcOpStart, SvcOcCfgNoRequiredParameter, cfg.Protocol, cfg.Transport, cfg.Entry, paramName)
			return false
		}
		value, _ := cfg.Params[paramName]
		log.Report(log.SrcSVC, SvcOpStart, SvcOcCfgInvalidParameterValue, cfg.Protocol, cfg.Transport, cfg.Entry, paramName, value)
		return false
	}

	return true
}
