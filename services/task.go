package services

import (
	"errors"

	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/defs"
	"github.com/stas-makutin/howeve/log"
	"github.com/stas-makutin/howeve/tasks"
)

// Task struct
type Task struct {
	cfg *config.Config
}

// NewTask func
func NewTask() *Task {
	t := &Task{}
	config.AddReader(t.readConfig)
	config.AddWriter(t.writeConfig)
	return t
}

func (t *Task) readConfig(cfg *config.Config, cfgError config.Error) {
	t.cfg = cfg
}

func (t *Task) writeConfig(cfg *config.Config) {
}

// Open func
func (t *Task) Open(ctx *tasks.ServiceTaskContext) error {

	for _, scfg := range t.cfg.Services {
		protocol, ok := ProtocolByName(scfg.Protocol)
		if !ok {
			log.Report(log.SrcSVC, SvcOpStart, SvcOcCfgUnknownProtocol, scfg.Protocol, scfg.Transport, scfg.Entry)
			continue
		}
		transport, ok := TransportByName(scfg.Transport)
		if !ok {
			log.Report(log.SrcSVC, SvcOpStart, SvcOcCfgUnknownTransport, scfg.Protocol, scfg.Transport, scfg.Entry)
			continue
		}
		pi, ok := Protocols[protocol]
		if !ok {
			log.Report(log.SrcSVC, SvcOpStart, SvcOcCfgProtocolNotSupported, scfg.Protocol, scfg.Transport, scfg.Entry)
			continue
		}
		pto, ok := pi.Transports[transport]
		if !ok {
			log.Report(log.SrcSVC, SvcOpStart, SvcOcCfgTransportNotSupported, scfg.Protocol, scfg.Transport, scfg.Entry)
			continue
		}
		ti, ok := Transports[transport]
		if !ok {
			log.Report(log.SrcSVC, SvcOpStart, SvcOcCfgTransportNotSupported, scfg.Protocol, scfg.Transport, scfg.Entry)
			continue
		}
		_ /*params*/, paramName, err := pto.Params.Merge(ti.Params).ParseAll(scfg.Params)
		if err != nil {
			if errors.Is(err, defs.ErrUnknownParamName) {
				log.Report(log.SrcSVC, SvcOpStart, SvcOcCfgUnknownParameter, scfg.Protocol, scfg.Transport, scfg.Entry, paramName)
				continue
			} else if errors.Is(err, defs.ErrNoRequiredParam) {
				log.Report(log.SrcSVC, SvcOpStart, SvcOcCfgNoRequiredParameter, scfg.Protocol, scfg.Transport, scfg.Entry, paramName)
				continue
			}
			value, _ := scfg.Params[paramName]
			log.Report(log.SrcSVC, SvcOpStart, SvcOcCfgInvalidParameterValue, scfg.Protocol, scfg.Transport, scfg.Entry, paramName, value)
			continue
		}
	}

	return nil
}

// Close func
func (t *Task) Close(ctx *tasks.ServiceTaskContext) error {
	return nil
}

// Stop func
func (t *Task) Stop(ctx *tasks.ServiceTaskContext) {
}
