package handlers

import (
	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/api"
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
	*api.Config
}

// ProtocolList - get list of supported protocols request
type ProtocolList struct {
	RequestHeader
}

// ProtocolListResult - get list of supported protocols response
type ProtocolListResult struct {
	ResponseHeader
	*api.ProtocolListResult
}

// TransportList - get list of all available transports despite the protocol
type TransportList struct {
	RequestHeader
}

// TransportListResult - the response to get list of all available transports despite the protocol
type TransportListResult struct {
	ResponseHeader
	*api.TransportListResult
}

// ProtocolInfo - get protocol(s) detailed information
type ProtocolInfo struct {
	RequestHeader
	*api.ProtocolInfo
}

// ProtocolInfoResult - the response to get list of all available transports despite the protocol
type ProtocolInfoResult struct {
	ResponseHeader
	*api.ProtocolInfoResult
}

// AddService - add new service
type AddService struct {
	RequestHeader
	*api.ServiceEntry
}

// AddServiceResult - add new service result
type AddServiceResult struct {
	ResponseHeader
	*api.StatusReply
}

// RemoveService - remove service request
type RemoveService struct {
	RequestHeader
	*api.ServiceID
}

// RemoveServiceResult - remove service result
type RemoveServiceResult struct {
	ResponseHeader
	*api.StatusReply
}

// ChangeServiceAlias - change service alias request
type ChangeServiceAlias struct {
	RequestHeader
	*api.ChangeServiceAlias
}

// ChangeServiceAlias - change service alias result
type ChangeServiceAliasResult struct {
	ResponseHeader
	*api.StatusReply
}

// ServiceStatus - get service status
type ServiceStatus struct {
	RequestHeader
	*api.ServiceID
}

// ServiceStatusResult - get service status
type ServiceStatusResult struct {
	ResponseHeader
	*api.StatusReply
}

// ListServices - get list of services request
type ListServices struct {
	RequestHeader
	*api.ListServices
}

// ListServicesResult - get list of services result envelope
type ListServicesResult struct {
	ResponseHeader
	*api.ListServicesResult
}

// SendToService - send message to service
type SendToService struct {
	RequestHeader
	*api.SendToService
}

// SendToServiceResult - send message to service result
type SendToServiceResult struct {
	ResponseHeader
	*api.SendToServiceResult
}

// ProtocolDiscoveryStarted event contains information about started discovery query
type ProtocolDiscoveryStarted struct {
	Header
	*api.ProtocolDiscoveryStarted
}

// ProtocolDiscoveryFinished event contains discovery query results
type ProtocolDiscoveryFinished struct {
	Header
	*api.ProtocolDiscoveryResult
}

// ProtocolDiscover defines protocol discover query event
type ProtocolDiscover struct {
	RequestHeader
	*api.ProtocolDiscover
}

// ProtocolDiscoverResult - discover query results
type ProtocolDiscoverResult struct {
	ResponseHeader
	*api.ProtocolDiscoverResult
}

// ProtocolDiscovery defines protocol discovery event
type ProtocolDiscovery struct {
	RequestHeader
	*api.ProtocolDiscovery
}

// ProtocolDiscoveryResult defines protocol discovery result event
type ProtocolDiscoveryResult struct {
	ResponseHeader
	*api.ProtocolDiscoveryResult
}

// NewMessage event contains information about new message
type NewMessage struct {
	Header
	*api.MessageEntry
}

// DropMessage event sent when a message gets removed from message log
type DropMessage struct {
	Header
	*api.MessageEntry
}

// UpdateMessageState event notifies about message state change
type UpdateMessageState struct {
	Header
	*api.UpdateMessageState
}

// GetMessage - get message request
type GetMessage struct {
	RequestHeader
	ID uuid.UUID
}

// GetMessageResult - get message result
type GetMessageResult struct {
	ResponseHeader
	*api.MessageEntry
}

// ListMessages defines list messages request
type ListMessages struct {
	RequestHeader
	*api.ListMessages
}

// ListMessagesResult defines list messages result
type ListMessagesResult struct {
	ResponseHeader
	*api.ListMessagesResult
}
