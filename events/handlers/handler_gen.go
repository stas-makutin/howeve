package handlers

import (
	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/defs"
	"github.com/stas-makutin/howeve/tasks"
)

func handleRestart(event *Restart) {
	Dispatcher.Send(&RestartResult{ResponseHeader: event.Associate()})
	go tasks.StopServiceTasks()
}

func handleConfigGet(event *ConfigGet, cfg *api.Config) {
	Dispatcher.Send(&ConfigGetResult{Config: cfg, ResponseHeader: event.Associate()})
}

func handleProtocolList(event *ProtocolList) {
	r := &ProtocolListResult{ResponseHeader: event.Associate(), ProtocolListResult: &api.ProtocolListResult{}}
	for k, v := range defs.Protocols {
		r.Protocols = append(r.Protocols, &api.ProtocolListEntry{ID: k, Name: v.Name})
	}
	Dispatcher.Send(r)
}

func handleTransportList(event *TransportList) {
	r := &TransportListResult{ResponseHeader: event.Associate(), TransportListResult: &api.TransportListResult{}}
	for k, v := range defs.Transports {
		r.Transports = append(r.Transports, &api.TransportListEntry{ID: k, Name: v.Name})
	}
	Dispatcher.Send(r)
}

func handleProtocolInfo(event *ProtocolInfo) {
	r := &ProtocolInfoResult{ResponseHeader: event.Associate()}

	tf := func(tid api.TransportIdentifier, pi *defs.ProtocolInfo, pie *api.ProtocolInfoEntry) {
		ptie := &api.ProtocolTransportInfoEntry{ID: tid, Valid: false}
		if ti, ok := defs.Transports[tid]; ok {
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
	pf := func(pid api.ProtocolIdentifier, pi *defs.ProtocolInfo) {
		pie := &api.ProtocolInfoEntry{ID: pid, Valid: true, Name: pi.Name}
		if event.ProtocolInfo != nil && len(event.Transports) > 0 {
			for _, t := range event.Transports {
				tf(t, pi, pie)
			}
		} else {
			for tid := range pi.Transports {
				tf(tid, pi, pie)
			}
		}
		if r.ProtocolInfoResult == nil {
			r.ProtocolInfoResult = &api.ProtocolInfoResult{}
		}
		r.Protocols = append(r.Protocols, pie)
	}

	if event.ProtocolInfo != nil && len(event.Protocols) > 0 {
		for _, pid := range event.Protocols {
			if pi, ok := defs.Protocols[pid]; ok {
				pf(pid, pi)
			} else {
				r.Protocols = append(r.Protocols, &api.ProtocolInfoEntry{ID: pid, Valid: false})
			}
		}
	} else {
		for pid, pi := range defs.Protocols {
			pf(pid, pi)
		}
	}

	Dispatcher.Send(r)
}
