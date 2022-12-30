package core

import (
	"fmt"

	"github.com/stas-makutin/howeve/api"
)

type ServiceEntryData struct {
	api.ServiceKey
	Alias  string
	Params Parameters
}

func NewServiceEntryData(protocols []*api.ProtocolInfoEntry) *ServiceEntryData {
	result := &ServiceEntryData{}
	if len(protocols) > 0 {
		protocol := protocols[0]
		result.Protocol = protocol.ID
		if len(protocol.Transports) > 0 {
			result.Transport = protocol.Transports[0].ID
		}
	}
	return result
}

func (sed *ServiceEntryData) ChangeProtocol(protocols []*api.ProtocolInfoEntry, index int) bool {
	if index < 0 || index >= len(protocols) {
		return false
	}
	protocol := protocols[index]
	if sed.Protocol == protocol.ID {
		return false
	}

	transportNotSupported := true
	for _, t := range protocol.Transports {
		if t.ID == sed.Transport {
			transportNotSupported = false
			break
		}
	}
	sed.Protocol = protocol.ID
	if transportNotSupported && len(protocol.Transports) > 0 {
		sed.Transport = protocol.Transports[0].ID
	}
	return true
}

func (sed *ServiceEntryData) ChangeTransport(protocols []*api.ProtocolInfoEntry, index int) bool {
	if index < 0 || len(protocols) == 0 {
		return false
	}

	var transports []*api.ProtocolTransportInfoEntry
	for _, p := range protocols {
		if p.ID == sed.Protocol {
			transports = p.Transports
		}
	}
	if index >= len(transports) {
		return false
	}

	transportID := transports[index].ID
	if transportID == sed.Transport {
		return false
	}

	sed.Transport = transportID
	return true
}

func (sed *ServiceEntryData) ValidateEntryAndAlias(services []api.ListServicesEntry, protocols []*api.ProtocolInfoEntry) (entryMessage, aliasMessage string) {
	for _, s := range services {
		if s.Alias != "" && s.Alias == sed.Alias {
			aliasMessage = fmt.Sprintf("Service with alias '%s' already exists", s.Alias)
		}
		if s.Protocol == sed.Protocol && s.Transport == sed.Transport && s.Entry == sed.Entry {
			protocol, transport := sed.ProtocolAndTransport(protocols)
			if protocol == nil || transport == nil {
				entryMessage = "Data integrity error - try to refresh the page"
			} else {
				entryMessage = fmt.Sprintf(
					"Service with '%s' entry already exists for %s (%v) protocol and %s (%v) transport",
					s.Entry,
					protocol.Name, s.Protocol,
					transport.Name, s.Transport,
				)
			}
		}
		if aliasMessage != "" && entryMessage != "" {
			break
		}
	}
	return
}

func (sed *ServiceEntryData) ProtocolAndTransport(protocols []*api.ProtocolInfoEntry) (protocol *api.ProtocolInfoEntry, transport *api.ProtocolTransportInfoEntry) {
	for _, p := range protocols {
		if p.ID == sed.Protocol {
			protocol = p
			for _, t := range protocol.Transports {
				if t.ID == sed.Transport {
					transport = t
					break
				}
			}
			break
		}
	}
	return
}
