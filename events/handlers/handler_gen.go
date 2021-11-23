package handlers

import (
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/defs"
	"github.com/stas-makutin/howeve/services"
	"github.com/stas-makutin/howeve/tasks"
)

func handleRestart(event *Restart) {
	Dispatcher.Send(&RestartResult{ResponseHeader: event.Associate()})
	go tasks.StopServiceTasks()
}

func handleConfigGet(event *ConfigGet, cfg *config.Config) {
	Dispatcher.Send(&ConfigGetResult{Config: *cfg, ResponseHeader: event.Associate()})
}

func handleProtocolList(event *ProtocolList) {
	r := &ProtocolListResult{ResponseHeader: event.Associate()}
	for k, v := range services.Protocols {
		r.Protocols = append(r.Protocols, &ProtocolListEntry{ID: uint8(k), Name: v.Name})
	}
	Dispatcher.Send(r)
}

func handleTransportList(event *TransportList) {
	r := &TransportListResult{ResponseHeader: event.Associate()}
	for k, v := range services.Transports {
		r.Transports = append(r.Transports, &TransportListEntry{ID: uint8(k), Name: v.Name})
	}
	Dispatcher.Send(r)
}

func handleProtocolInfo(event *ProtocolInfo) {
	r := &ProtocolInfoResult{ResponseHeader: event.Associate()}

	tf := func(tid defs.TransportIdentifier, pi *defs.ProtocolInfo, pie *ProtocolInfoEntry) {
		ptie := &ProtocolTransportInfoEntry{ID: tid, Valid: false}
		if ti, ok := services.Transports[tid]; ok {
			ptie.Name = ti.Name
			if pto, ok := pi.Transports[tid]; ok {
				ptie.Valid = true
				ptie.Discoverable = pto.DiscoveryFunc != nil
				ptie.Params = buildParamsInfo(pto.Params.Merge(ti.Params))
				ptie.DiscoveryParams = buildParamsInfo(pto.DiscoveryParams)
			}
		}
		pie.Transports = append(pie.Transports, ptie)
	}
	pf := func(pid defs.ProtocolIdentifier, pi *defs.ProtocolInfo) {
		pie := &ProtocolInfoEntry{ID: pid, Valid: true, Name: pi.Name}
		if event.Filter != nil && len(event.Filter.Transports) > 0 {
			for _, t := range event.Filter.Transports {
				tf(defs.TransportIdentifier(t), pi, pie)
			}
		} else {
			for tid := range pi.Transports {
				tf(tid, pi, pie)
			}
		}
		r.Protocols = append(r.Protocols, pie)
	}

	if event.Filter != nil && len(event.Filter.Protocols) > 0 {
		for _, p := range event.Filter.Protocols {
			pid := defs.ProtocolIdentifier(p)
			if pi, ok := services.Protocols[pid]; ok {
				pf(pid, pi)
			} else {
				r.Protocols = append(r.Protocols, &ProtocolInfoEntry{ID: p, Valid: false})
			}
		}
	} else {
		for pid, pi := range services.Protocols {
			pf(pid, pi)
		}
	}

	Dispatcher.Send(r)
}

func handleProtocolDiscovery(event *ProtocolDiscovery) {
	r := &ProtocolDiscoveryResult{ResponseHeader: event.Associate(), ProtocolDiscoveryQueryResult: &ProtocolDiscoveryQueryResult{}}
	if pat, ei := findProtocolAndTransport(event.Protocol, event.Transport); ei == nil {
		if pat.options.DiscoveryFunc != nil {
			if params, errorInfo := event.Params.parse(pat.options.DiscoveryParams); errorInfo != nil {
				r.ProtocolDiscoveryQueryResult.Error = errorInfo
			} else {
				go func() {
					if serviceEntries, err := pat.options.DiscoveryFunc(event.Context(), params); err == nil {
						if len(serviceEntries) > 0 {
							r.Services = make([]*ServiceEntryDetails, 0, len(serviceEntries))
							for _, serviceEntry := range serviceEntries {
								r.Services = append(r.Services, &ServiceEntryDetails{
									ServiceEntry: ServiceEntry{
										ServiceKey: ServiceKey{
											Protocol:  serviceEntry.Key.Protocol,
											Transport: serviceEntry.Key.Transport,
											Entry:     serviceEntry.Key.Entry,
										},
										Params: NewParamsValues(serviceEntry.Params),
									},
									Description: serviceEntry.Description,
								})
							}
						}
					} else {
						// TODO error reporting
					}
					Dispatcher.Send(r)
				}()
				return
			}
		} else {
			r.ProtocolDiscoveryQueryResult.Error = NewErrorInfo(
				ErrorNoDiscovery, pat.protocol.Name, event.Protocol, pat.transport.Name, event.Transport,
			)
		}
	} else {
		r.ProtocolDiscoveryQueryResult.Error = ei
	}

	Dispatcher.Send(r)
}
