package core

import (
	"fmt"
	"strconv"

	"github.com/stas-makutin/howeve/api"
)

type ProtocolsInfoWrapper struct {
	Info         *api.ProtocolInfoEntry
	transportIDs map[api.TransportIdentifier]*api.ProtocolTransportInfoEntry
}

type ProtocolsWrapper struct {
	Protocols   []*api.ProtocolInfoEntry
	protocolIDs map[api.ProtocolIdentifier]*ProtocolsInfoWrapper
}

func NewProtocolsWrapper(protocols *api.ProtocolInfoResult) *ProtocolsWrapper {
	result := &ProtocolsWrapper{}
	if protocols != nil {
		result.Protocols = protocols.Protocols
		result.protocolIDs = make(map[api.ProtocolIdentifier]*ProtocolsInfoWrapper)
		for _, protocol := range protocols.Protocols {
			p := &ProtocolsInfoWrapper{Info: protocol, transportIDs: make(map[api.TransportIdentifier]*api.ProtocolTransportInfoEntry)}
			result.protocolIDs[protocol.ID] = p
			for _, transport := range protocol.Transports {
				p.transportIDs[transport.ID] = transport
			}
		}
	}
	return result
}

func (pw *ProtocolsWrapper) ProtocolAndTransport(protocolID api.ProtocolIdentifier, transportID api.TransportIdentifier) (protocol *api.ProtocolInfoEntry, transport *api.ProtocolTransportInfoEntry) {
	if p, ok := pw.protocolIDs[protocolID]; ok {
		protocol = p.Info
		if t, ok := p.transportIDs[transportID]; ok {
			transport = t
		}
	}
	return
}

func (pw *ProtocolsWrapper) ProtocolAndTransportNames(protocolID api.ProtocolIdentifier, transportID api.TransportIdentifier) (protocolName, transportName string) {
	pi, ti := pw.ProtocolAndTransport(protocolID, transportID)
	if pi != nil {
		protocolName = pi.Name
	}
	if ti != nil {
		transportName = ti.Name
	}
	return
}

func (pw *ProtocolsWrapper) ProtocolAndTransportFullNames(protocolID api.ProtocolIdentifier, transportID api.TransportIdentifier) (protocolName, transportName string) {
	protocolName = strconv.Itoa(int(protocolID))
	transportName = strconv.Itoa(int(transportID))
	pi, ti := pw.ProtocolAndTransport(protocolID, transportID)
	if pi != nil && pi.Name != "" {
		protocolName = pi.Name + " (" + protocolName + ")"
	}
	if ti != nil && ti.Name != "" {
		transportName = ti.Name + " (" + transportName + ")"
	}
	return
}

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

func (sed *ServiceEntryData) ChangeTransport(protocols *ProtocolsWrapper, index int) bool {
	if index < 0 || len(protocols.Protocols) == 0 {
		return false
	}

	piw, ok := protocols.protocolIDs[sed.Protocol]
	if !ok {
		return false
	}
	if index >= len(piw.Info.Transports) {
		return false
	}

	transportID := piw.Info.Transports[index].ID
	if transportID == sed.Transport {
		return false
	}

	sed.Transport = transportID
	return true
}

func (sed *ServiceEntryData) ValidateEntryAndAlias(services []api.ListServicesEntry, protocols *ProtocolsWrapper) (entryMessage, aliasMessage string) {
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

func (sed *ServiceEntryData) ProtocolAndTransport(protocols *ProtocolsWrapper) (*api.ProtocolInfoEntry, *api.ProtocolTransportInfoEntry) {
	return protocols.ProtocolAndTransport(sed.Protocol, sed.Transport)
}
