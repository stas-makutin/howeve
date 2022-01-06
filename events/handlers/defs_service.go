package handlers

import (
	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/defs"
)

// ServiceID
type ServiceID struct {
	*defs.ServiceKey
	Alias string `json:"alias,omitempty"`
}

// ServiceEntry declares the service
type ServiceEntry struct {
	*defs.ServiceKey
	Params defs.RawParamValues `json:"params,omitempty"`
	Alias  string              `json:"alias,omitempty"`
}

// AddService - add new service
type AddService struct {
	RequestHeader
	*ServiceEntry
}

// ServiceReply - add new service reply
type ServiceReply struct {
	Error   *ErrorInfo `json:"error,omitempty"`
	Success bool       `json:"success,omitempty"`
}

// AddServiceResult - add new service result
type AddServiceResult struct {
	ResponseHeader
	*StatusReply
}

// RemoveService - remove service request
type RemoveService struct {
	RequestHeader
	*ServiceID
}

// RemoveServiceResult - remove service result
type RemoveServiceResult struct {
	ResponseHeader
	*StatusReply
}

// ChangeServiceAliasQuery - change service alias request payload
type ChangeServiceAliasQuery struct {
	*ServiceID
	NewAlias string `json:"newAlias,omitempty"`
}

// ChangeServiceAlias - change service alias request
type ChangeServiceAlias struct {
	RequestHeader
	*ChangeServiceAliasQuery
}

// ChangeServiceAlias - change service alias result
type ChangeServiceAliasResult struct {
	ResponseHeader
	*StatusReply
}

// ServiceStatus - get service status
type ServiceStatus struct {
	RequestHeader
	*ServiceID
}

// ServiceStatusResult - get service status
type ServiceStatusResult struct {
	ResponseHeader
	*StatusReply
}

// ListServicesInput - get list of services request inputs
type ListServicesInput struct {
	Protocols  []defs.ProtocolIdentifier  `json:"protocols,omitempty"`
	Transports []defs.TransportIdentifier `json:"transports,omitempty"`
	Entries    []string                   `json:"entries,omitempty"`
	Aliases    []string                   `json:"aliases,omitempty"`
}

// ListServices - get list of services request
type ListServices struct {
	RequestHeader
	*ListServicesInput
}

// ListServicesEntry - service information for services list result
type ListServicesEntry struct {
	*ServiceID
	*StatusReply
}

// ListServicesOutput - get list of services result
type ListServicesOutput struct {
	Services []ListServicesEntry `json:"services,omitempty"`
}

// ListServicesResult - get list of services result envelope
type ListServicesResult struct {
	ResponseHeader
	*ListServicesOutput
}

// SendToService - send message to service request input
type SendToServiceInput struct {
	*ServiceID
	Payload []byte `json:"payload,omitempty"`
}

// SendToService - send message to service
type SendToService struct {
	RequestHeader
	*SendToServiceInput
}

// SendToServiceOutput - send message to service result payload
type SendToServiceOutput struct {
	ResponseHeader
	*StatusReply
	*defs.Message
}

// SendToServiceResult - send message to service result
type SendToServiceResult struct {
	ResponseHeader
	*SendToServiceOutput
}

// DiscoveryStarted event contains information about started discovery query
type DiscoveryStarted struct {
	Header
	ID        uuid.UUID                `json:"id"`
	Protocol  defs.ProtocolIdentifier  `json:"protocol"`
	Transport defs.TransportIdentifier `json:"transport"`
	Params    defs.RawParamValues      `json:"params,omitempty"`
}

// DiscoveryResult defines discovery results
type DiscoveryResult struct {
	ID      uuid.UUID              `json:"id"`
	Entries []*defs.DiscoveryEntry `json:"entries"`
	Error   *ErrorInfo             `json:"error,omitempty"`
}

// DiscoveryFinished event contains discovery query results
type DiscoveryFinished struct {
	Header
	*DiscoveryResult
}

// ProtocolDiscoverInput - discover query input parameters
type ProtocolDiscoverInput struct {
	Protocol  defs.ProtocolIdentifier  `json:"protocol"`
	Transport defs.TransportIdentifier `json:"transport"`
	Params    defs.RawParamValues      `json:"params,omitempty"`
}

// ProtocolDiscover defines protocol discover query event
type ProtocolDiscover struct {
	RequestHeader
	*ProtocolDiscoverInput
}

// ProtocolDiscoverInput - discover query output parameters
type ProtocolDiscoverOutput struct {
	ID    *uuid.UUID `json:"id,omitempty"`
	Error *ErrorInfo `json:"error,omitempty"`
}

// ProtocolDiscoverResult - discover query results
type ProtocolDiscoverResult struct {
	ResponseHeader
	*ProtocolDiscoverOutput
}

// ProtocolDiscoveryInput defines input of protocol discovery event
type ProtocolDiscoveryInput struct {
	ID   uuid.UUID `json:"id"`
	Stop bool      `json:"stop"`
}

// ProtocolDiscovery defines protocol discovery event
type ProtocolDiscovery struct {
	RequestHeader
	*ProtocolDiscoveryInput
}

// ProtocolDiscoveryResult defines protocol discovery result event
type ProtocolDiscoveryResult struct {
	ResponseHeader
	*DiscoveryResult
}
