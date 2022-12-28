package core

import "github.com/stas-makutin/howeve/api"

func ProtocolAndTransportName(protocol api.ProtocolIdentifier, transport api.TransportIdentifier, protocols *api.ProtocolInfoResult) (protocolName, transportName string) {
	if protocols != nil && len(protocols.Protocols) > 0 {
		for _, p := range protocols.Protocols {
			if p.ID == protocol {
				protocolName = p.Name
				for _, t := range p.Transports {
					if t.ID == transport {
						transportName = t.Name
						return
					}
				}
			}
		}
	}
	return
}
