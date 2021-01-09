package handlers

import (
	"errors"
	"fmt"

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

func buildParamsInfo(p defs.Params) (pie map[string]*ParamInfoEntry) {
	pie = make(map[string]*ParamInfoEntry)
	for name, pi := range p {
		if pi.Flags&defs.ParamFlagConst == 0 {
			pie[name] = &ParamInfoEntry{
				Description:  pi.Description,
				Type:         pi.Type.String(),
				DefaultValue: pi.DefaultValue,
				EnumValues:   pi.EnumValues,
			}
		}
	}
	return
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

// ParamsValues type
type ParamsValues map[string]string

// NewParamsValues creates reporting parameter values from service parameter values
func NewParamsValues(pv defs.ParamValues) (r ParamsValues) {
	if len(pv) > 0 {
		r = make(ParamsValues)
		for name, value := range pv {
			r[name] = fmt.Sprint(value)
		}
	}
	return
}

// Parse function parses input parameters according their definition
func (pv ParamsValues) Parse(p defs.Params) (defs.ParamValues, *ErrorInfo) {
	values, name, err := p.ParseAll(pv)
	if err != nil {
		if errors.Is(err, defs.ErrUnknownParamName) {
			return nil, NewErrorInfo(ErrorUnknownParameter, name)
		} else if errors.Is(err, defs.ErrNoRequiredParam) {
			return nil, NewErrorInfo(ErrorNoRequiredParameter, name)
		}
		value, _ := pv[name]
		return nil, NewErrorInfo(ErrorInvalidParameterValue, value, name)
	}
	return values, nil
}

// ProtocolDiscoveryQuery - discovery input parameters
type ProtocolDiscoveryQuery struct {
	Protocol  defs.ProtocolIdentifier  `json:"protocol"`
	Transport defs.TransportIdentifier `json:"transport"`
	Params    ParamsValues             `json:"params,omitempty"`
}

// ServiceEntry - service entry description
type ServiceEntry struct {
	Protocol  defs.ProtocolIdentifier  `json:"protocol"`
	Transport defs.TransportIdentifier `json:"transport"`
	Entry     string                   `json:"entry"`
	Params    ParamsValues             `json:"params,omitempty"`
}

// ServiceEntryDetails - service entry with details
type ServiceEntryDetails struct {
	ServiceEntry
	Description string `json:"description,omitempty"`
}

// ProtocolDiscovery - discovery available services of protocol using specific transport
type ProtocolDiscovery struct {
	RequestHeader
	*ProtocolDiscoveryQuery
}

// ProtocolDiscoveryQueryResult - discovery query results
type ProtocolDiscoveryQueryResult struct {
	Error    *ErrorInfo             `json:"error,omitempty"`
	Services []*ServiceEntryDetails `json:"services,omitempty"`
}

// ProtocolDiscoveryResult - discovery results
type ProtocolDiscoveryResult struct {
	ResponseHeader
	*ProtocolDiscoveryQueryResult
}

// ServiceEntryWithAlias - service entry with alias
type ServiceEntryWithAlias struct {
	ServiceEntry
	Alias string `json:"alias,omitempty"`
}

// AddService - add new service
type AddService struct {
	RequestHeader
	*ServiceEntryWithAlias
}

// AddServiceReply - add new service reply
type AddServiceReply struct {
	Error   *ErrorInfo `json:"error,omitempty"`
	Success bool       `json:"success,omitempty"`
}

// AddServiceResult - add new service result
type AddServiceResult struct {
	ResponseHeader
	*AddServiceReply
}
