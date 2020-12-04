package handlers

import (
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/services"
)

// Restart - restart the service
type Restart struct {
	RequestHeader
}

// RestartResult - restart the service result
type RestartResult struct {
	ResponseHeader
}

// ConfigGet - get config event
type ConfigGet struct {
	RequestHeader
}

// ConfigGetResult - config data event
type ConfigGetResult struct {
	ResponseHeader
	config.Config
}

// ProtocolListEntry - list of supported protocols
type ProtocolListEntry struct {
	ID   uint8  `json:"id"`
	Name string `json:"name"`
}

// ProtocolList - get list of supported protocols request
type ProtocolList struct {
	RequestHeader
}

// ProtocolListResult - get list of supported protocols response
type ProtocolListResult struct {
	ResponseHeader
	Protocols []*ProtocolListEntry
}

// TransportListEntry - list of all available transports despite the protocol
type TransportListEntry struct {
	ID   uint8  `json:"id"`
	Name string `json:"name"`
}

// TransportList - get list of all available transports despite the protocol
type TransportList struct {
	RequestHeader
}

// TransportListResult - the response to get list of all available transports despite the protocol
type TransportListResult struct {
	ResponseHeader
	Transports []*TransportListEntry
}

// ParamInfoEntry - information about named parameter
type ParamInfoEntry struct {
	Description  string   `json:"description"`
	Type         string   `json:"type"`
	DefaultValue string   `json:"defaultValue,omitempty"`
	EnumValues   []string `json:"enumValues,omitempty"`
}

func buildParamsInfo(p services.Params) (pie map[string]*ParamInfoEntry) {
	pie = make(map[string]*ParamInfoEntry)
	for name, pi := range p {
		pie[name] = &ParamInfoEntry{
			Description:  pi.Description,
			Type:         pi.Type.String(),
			DefaultValue: pi.DefaultValue,
			EnumValues:   pi.EnumValues,
		}
	}
	return
}

// ProtocolTransportInfoEntry - protocol detaild information
type ProtocolTransportInfoEntry struct {
	ID              services.TransportIdentifier `json:"id"`
	Valid           bool                         `json:"valid"`
	Name            string                       `json:"name,omitempty"`
	Params          map[string]*ParamInfoEntry   `json:"params,omitempty"`
	Discoverable    bool                         `json:"discoverable"`
	DiscoveryParams map[string]*ParamInfoEntry   `json:"discoveryParams,omitempty"`
}

// ProtocolInfoEntry - protocol detaild information
type ProtocolInfoEntry struct {
	ID         services.ProtocolIdentifier   `json:"id"`
	Valid      bool                          `json:"valid"`
	Name       string                        `json:"name,omitempty"`
	Transports []*ProtocolTransportInfoEntry `json:"transports,omitempty"`
}

// ProtocolInfoFilter - protocols/transport filter
type ProtocolInfoFilter struct {
	Protocols  []services.ProtocolIdentifier  `json:"protocols,omitempty"`
	Transports []services.TransportIdentifier `json:"transports,omitempty"`
}

// ProtocolInfo - get protocol(s) detailed information
type ProtocolInfo struct {
	RequestHeader
	Filter *ProtocolInfoFilter
}

// ProtocolInfoResult - the response to get list of all available transports despite the protocol
type ProtocolInfoResult struct {
	ResponseHeader
	Protocols []*ProtocolInfoEntry
}

// ProtocolDiscoveryQuery - discovery input parameters
type ProtocolDiscoveryQuery struct {
	Protocol  services.ProtocolIdentifier  `json:"protocol"`
	Transport services.TransportIdentifier `json:"transport"`
	Params    map[string]string            `json:"params,omitempty"`
}

// ServiceEntry - service entry description
type ServiceEntry struct {
	Protocol  services.ProtocolIdentifier  `json:"protocol"`
	Transport services.TransportIdentifier `json:"transport"`
	Entry     string                       `json:"entry"`
	Params    map[string]string            `json:"params,omitempty"`
}

// ProtocolDiscovery - discovery available services of protocol using specific transport
type ProtocolDiscovery struct {
	RequestHeader
	*ProtocolDiscoveryQuery
}

// ProtocolDiscoveryQueryResult - discovery query results
type ProtocolDiscoveryQueryResult struct {
	Valid    bool            `json:"valid,omitempty"`
	Services []*ServiceEntry `json:"services,omitempty"`
}

// ProtocolDiscoveryResult - discovery results
type ProtocolDiscoveryResult struct {
	ResponseHeader
	*ProtocolDiscoveryQueryResult
}
