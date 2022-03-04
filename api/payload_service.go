package api

import (
	"github.com/google/uuid"
)

// ServiceKey defines service unique identifier/key
type ServiceKey struct {
	Protocol  ProtocolIdentifier  `json:"protocol"`
	Transport TransportIdentifier `json:"transport"`
	Entry     string              `json:"entry"`
}

// ServiceID defines service identification in API
type ServiceID struct {
	*ServiceKey
	Alias string `json:"alias,omitempty"`
}

// ServiceEntry declares service identification and configuration
type ServiceEntry struct {
	*ServiceKey
	Params RawParamValues `json:"params,omitempty"`
	Alias  string         `json:"alias,omitempty"`
}

// ChangeServiceAlias - change service alias request payload
type ChangeServiceAlias struct {
	*ServiceID
	NewAlias string `json:"newAlias,omitempty"`
}

// ListServices - get list of services request payload
type ListServices struct {
	Protocols  []ProtocolIdentifier  `json:"protocols,omitempty"`
	Transports []TransportIdentifier `json:"transports,omitempty"`
	Entries    []string              `json:"entries,omitempty"`
	Aliases    []string              `json:"aliases,omitempty"`
}

// ListServicesEntry - service information for services list result
type ListServicesEntry struct {
	*ServiceID
	*StatusReply
}

// ListServicesResult - get list of services query result
type ListServicesResult struct {
	Services []ListServicesEntry `json:"services,omitempty"`
}

// SendToService - send message to service request payload
type SendToService struct {
	*ServiceID
	Payload []byte `json:"payload,omitempty"`
}

// SendToServiceResult - send message to service result payload
type SendToServiceResult struct {
	*StatusReply
	*Message
}

// ProtocolDiscover - discover query request payload
type ProtocolDiscover struct {
	Protocol  ProtocolIdentifier  `json:"protocol"`
	Transport TransportIdentifier `json:"transport"`
	Params    RawParamValues      `json:"params,omitempty"`
}

// ProtocolDiscoverResult - discover query response payload
type ProtocolDiscoverResult struct {
	ID    *uuid.UUID `json:"id,omitempty"`
	Error *ErrorInfo `json:"error,omitempty"`
}

// ProtocolDiscovery defines protocol discovery query payload
type ProtocolDiscovery struct {
	ID   uuid.UUID `json:"id"`
	Stop bool      `json:"stop"`
}

// DiscoveryEntry - information about service instance found during discovery
type DiscoveryEntry struct {
	ServiceKey
	ParamValues `json:"params,omitempty"`
	Description string `json:"description,omitempty"`
}

// ProtocolDiscoveryResult defines discovery query response payload
type ProtocolDiscoveryResult struct {
	ID      uuid.UUID         `json:"id"`
	Entries []*DiscoveryEntry `json:"entries"`
	Error   *ErrorInfo        `json:"error,omitempty"`
}

// DiscoveryRequest contains information about started discovery query
type ProtocolDiscoveryStarted struct {
	ProtocolDiscover
	ID uuid.UUID `json:"id"`
}
