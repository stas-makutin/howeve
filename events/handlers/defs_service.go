package handlers

import (
	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/defs"
)

// ServiceID
type ServiceID struct {
	*ServiceKey
	Alias string `json:"alias,omitempty"`
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

// SendToService - send message to service
type SendToService struct {
	RequestHeader
	*ServiceID
}

// SendToServiceResult - send message to service result
type SendToServiceResult struct {
	ResponseHeader
}

// DiscoveryStarted event contains information about started discovery query
type DiscoveryStarted struct {
	Header
	ID        uuid.UUID                `json:"id"`
	Protocol  defs.ProtocolIdentifier  `json:"protocol"`
	Transport defs.TransportIdentifier `json:"transport"`
	Params    defs.RawParamValues      `json:"params,omitempty"`
}

type DiscoveryEntry struct {
	ServiceKey
	Params      defs.ParamValues `json:"params,omitempty"`
	Description string           `json:"description,omitempty"`
}

// DiscoveryFinished event contains discovery query results
type DiscoveryFinished struct {
	Header
	ID      uuid.UUID        `json:"id"`
	Entries []DiscoveryEntry `json:"entries"`
	Error   *ErrorInfo       `json:"error,omitempty"`
}

// // ProtocolDiscovery - discovery available services of protocol using specific transport
// type ProtocolDiscovery struct {
// 	RequestHeader
// 	*ProtocolDiscoveryQuery
// }

// // ProtocolDiscoveryQueryResult - discovery query results
// type ProtocolDiscoveryQueryResult struct {
// 	Error    *ErrorInfo             `json:"error,omitempty"`
// 	Services []*ServiceEntryDetails `json:"services,omitempty"`
// }

// // ProtocolDiscoveryResult - discovery results
// type ProtocolDiscoveryResult struct {
// 	ResponseHeader
// 	*ProtocolDiscoveryQueryResult
// }
