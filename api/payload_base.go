package api

// ProtocolListEntry - list of supported protocols
type ProtocolListEntry struct {
	ID   ProtocolIdentifier `json:"id"`
	Name string             `json:"name"`
}

// ProtocolListResult - get list of supported protocols response
type ProtocolListResult struct {
	Protocols []*ProtocolListEntry
}

// TransportListEntry - list of all available transports despite the protocol
type TransportListEntry struct {
	ID   TransportIdentifier `json:"id"`
	Name string              `json:"name"`
}

// TransportListResult - the response to get list of all available transports despite the protocol
type TransportListResult struct {
	Transports []*TransportListEntry
}

// ParamInfoEntry - information about named parameter
type ParamInfoEntry struct {
	Description  string   `json:"description"`
	Type         string   `json:"type"`
	DefaultValue string   `json:"defaultValue,omitempty"`
	EnumValues   []string `json:"enumValues,omitempty"`
}

// RawParamValues defines raw parameter values - before parsing
type RawParamValues map[string]string

// ParamValues type defines named parameter values
type ParamValues map[string]interface{}

// ProtocolTransportInfoEntry - protocol detaild information
type ProtocolTransportInfoEntry struct {
	ID              TransportIdentifier        `json:"id"`
	Valid           bool                       `json:"valid"`
	Name            string                     `json:"name,omitempty"`
	Params          map[string]*ParamInfoEntry `json:"params,omitempty"`
	Discoverable    bool                       `json:"discoverable,omitempty"`
	DiscoveryParams map[string]*ParamInfoEntry `json:"discoveryParams,omitempty"`
}

// ProtocolInfoEntry - protocol detaild information
type ProtocolInfoEntry struct {
	ID         ProtocolIdentifier            `json:"id"`
	Valid      bool                          `json:"valid"`
	Name       string                        `json:"name,omitempty"`
	Transports []*ProtocolTransportInfoEntry `json:"transports,omitempty"`
}

// ProtocolInfo - get protocol(s) detailed information
type ProtocolInfo struct {
	Protocols  []ProtocolIdentifier  `json:"protocols,omitempty"`
	Transports []TransportIdentifier `json:"transports,omitempty"`
}

// ProtocolInfoResult - the response to get list of all available transports despite the protocol
type ProtocolInfoResult struct {
	Protocols []*ProtocolInfoEntry
}

// StatusReply - common operation status reply (success/error)
type StatusReply struct {
	Error   *ErrorInfo `json:"error,omitempty"`
	Success bool       `json:"success,omitempty"`
}

// PayloadMatch defines byte sequence to match with payload
// At == nil 	- payload must contain the content
// At >= 0 		- payload must include the content at provided position "At"
// At < 0 		- payload must include the content at provided position "len(payload) - len(content) + At + 1"
type PayloadMatch struct {
	Content []byte `json:"content,omitempty"`
	At      *int   `json:"at,omitempty"`
}
