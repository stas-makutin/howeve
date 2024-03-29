package api

import "fmt"

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

const (
	ParamTypeInt8   = "int8"
	ParamTypeInt16  = "int16"
	ParamTypeInt32  = "int32"
	ParamTypeInt64  = "int64"
	ParamTypeUint8  = "uint8"
	ParamTypeUint16 = "uint16"
	ParamTypeUint32 = "uint32"
	ParamTypeUint64 = "uint64"
	ParamTypeBool   = "bool"
	ParamTypeString = "string"
	ParamTypeEnum   = "enum"

	ParamTypeUint8Max = ^uint8(0)
	ParamTypeUint8Min = 0
	ParamTypeInt8Max  = int8(ParamTypeUint8Max >> 1)
	ParamTypeInt8Min  = -ParamTypeInt8Max - 1

	ParamTypeUint16Max = ^uint16(0)
	ParamTypeUint16Min = 0
	ParamTypeInt16Max  = int16(ParamTypeUint16Max >> 1)
	ParamTypeInt16Min  = -ParamTypeInt16Max - 1

	ParamTypeUint32Max = ^uint32(0)
	ParamTypeUint32Min = 0
	ParamTypeInt32Max  = int32(ParamTypeUint32Max >> 1)
	ParamTypeInt32Min  = -ParamTypeInt32Max - 1

	ParamTypeUint64Max = ^uint64(0)
	ParamTypeUint64Min = 0
	ParamTypeInt64Max  = int64(ParamTypeUint64Max >> 1)
	ParamTypeInt64Min  = -ParamTypeInt64Max - 1
)

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

// Raw converts parameter-values into their raw form (string-string)
func (pv ParamValues) Raw() (r RawParamValues) {
	if len(pv) > 0 {
		r = make(RawParamValues)
		for name, value := range pv {
			r[name] = fmt.Sprint(value)
		}
	}
	return
}

// Copy creates copy of parameter-values
func (pv ParamValues) Copy() ParamValues {
	rv := make(ParamValues)
	for k, v := range pv {
		rv[k] = v
	}
	return rv
}

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
