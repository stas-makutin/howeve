package handlers

import (
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/defs"
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

// ProtocolTransportInfoEntry - protocol detaild information
type ProtocolTransportInfoEntry struct {
	ID              defs.TransportIdentifier   `json:"id"`
	Valid           bool                       `json:"valid"`
	Name            string                     `json:"name,omitempty"`
	Params          map[string]*ParamInfoEntry `json:"params,omitempty"`
	Discoverable    bool                       `json:"discoverable,omitempty"`
	DiscoveryParams map[string]*ParamInfoEntry `json:"discoveryParams,omitempty"`
}

// ProtocolInfoEntry - protocol detaild information
type ProtocolInfoEntry struct {
	ID         defs.ProtocolIdentifier       `json:"id"`
	Valid      bool                          `json:"valid"`
	Name       string                        `json:"name,omitempty"`
	Transports []*ProtocolTransportInfoEntry `json:"transports,omitempty"`
}

// ProtocolInfoFilter - protocols/transport filter
type ProtocolInfoFilter struct {
	Protocols  []defs.ProtocolIdentifier  `json:"protocols,omitempty"`
	Transports []defs.TransportIdentifier `json:"transports,omitempty"`
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
	Protocol  defs.ProtocolIdentifier  `json:"protocol"`
	Transport defs.TransportIdentifier `json:"transport"`
	Params    defs.RawParamValues      `json:"params,omitempty"`
}

// ServiceKey - service identification/key
type ServiceKey struct {
	Protocol  defs.ProtocolIdentifier  `json:"protocol"`
	Transport defs.TransportIdentifier `json:"transport"`
	Entry     string                   `json:"entry"`
}

// ServiceEntry - service entry description
type ServiceEntry struct {
	ServiceKey
	Params defs.RawParamValues `json:"params,omitempty"`
}

// ServiceEntryDetails - service entry with details
type ServiceEntryDetails struct {
	ServiceEntry
	Description string `json:"description,omitempty"`
}
