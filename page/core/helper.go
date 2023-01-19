package core

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/stas-makutin/howeve/api"
)

var errorCodeToName map[api.ErrorCode]string = map[api.ErrorCode]string{
	api.ErrorUnknownProtocol:          "Unknown Protocol",
	api.ErrorUnknownTransport:         "Unknown Transport",
	api.ErrorInvalidProtocolTransport: "Invalid Protocol's Transport",
	api.ErrorUnknownParameter:         "Unknown Parameter",
	api.ErrorInvalidParameterValue:    "Invalid Parameter Value",
	api.ErrorNoRequiredParameter:      "No Required Parameter",
	api.ErrorNoDiscovery:              "No Discovery",
	api.ErrorDiscoveryBusy:            "Discovery Busy",
	api.ErrorNoDiscoveryID:            "No Discovery ID",
	api.ErrorDiscoveryPending:         "Discovery Pending",
	api.ErrorDiscoveryFailed:          "Discovery Failed",
	api.ErrorServiceNoKey:             "Service No Key",
	api.ErrorServiceNoID:              "Service No ID",
	api.ErrorServiceExists:            "Service Exists",
	api.ErrorServiceAliasExists:       "Service Alias Exists",
	api.ErrorServiceInitialize:        "Service Initialize",
	api.ErrorServiceKeyNotExists:      "Service Key Not Exists",
	api.ErrorServiceAliasNotExists:    "Service Alias Not Exists",
	api.ErrorServiceStatusBad:         "Service Status Bad",
	api.ErrorServiceBadPayload:        "Service Bad Payload",
	api.ErrorServiceSendBusy:          "Service Send Busy",
	api.ErrorOtherError:               "Other Error",
}

func ApiErrorName(code api.ErrorCode) string {
	name, ok := errorCodeToName[code]
	if ok {
		return name
	}
	return "Unknown"
}

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

func ArrangeServices(services []api.ListServicesEntry) {
	sort.SliceStable(services, func(i, j int) bool {
		if services[i].Protocol < services[j].Protocol {
			return true
		}
		if services[i].Protocol == services[j].Protocol {
			if services[i].Transport < services[j].Transport {
				return true
			}
			if services[i].Transport == services[j].Transport {
				if c := strings.Compare(services[i].Alias, services[j].Alias); c > 0 {
					return true
				} else if c == 0 {
					return strings.Compare(services[i].Entry, services[j].Entry) < 0
				}
			}
		}
		return false
	})
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
