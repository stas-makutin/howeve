package defs

import (
	"errors"

	"github.com/google/uuid"
)

// ServiceKey struct defines service unique identifier/key
type ServiceKey struct {
	Protocol  ProtocolIdentifier
	Transport TransportIdentifier
	Entry     string
}

// ServiceStatus describes the status of the service
type ServiceStatus struct {
}

// Service interface, defines minimal set of methods the service needs to support
type Service interface {
	Start()
	Stop()
	Status() ServiceStatus
	Send(payload []byte) (*Message, error)
}

// errors
var (
	// ErrServiceExists is the error in case if service already exists
	ErrServiceExists error = errors.New("the service already exists")
	// ErrAliasExists is the error in case if service already exists
	ErrAliasExists error = errors.New("the service alias already exists")
	// ErrProtocolNotSupported is the error in case if provided protocol is not supported
	ErrProtocolNotSupported error = errors.New("the protocol is not supported")
	// ErrTransportNotSupported is the error in case if provided transport is not supported for given protocol
	ErrTransportNotSupported error = errors.New("the transport is not supported")
	// ErrBadPayload returned by Send method in case if message's payload is not valid has no payload
	ErrBadPayload error = errors.New("the message's payload is not valid")
	// ErrSendBusy returned by Send method in case if service is unable to send message at this time
	ErrSendBusy error = errors.New("the service is too busy and unable to send the message")
)

// ParamNameOpenAttemptsInterval parameter name for the time interval between attempts to open serial port
const ParamNameOpenAttemptsInterval = "openAttemptsInterval"

// ServiceRegistry defines possible operations with services
type ServiceRegistry interface {
	Discover(protocol ProtocolIdentifier, transport TransportIdentifier, params RawParamValues) (uuid.UUID, error)
	Discovery(id uuid.UUID, stop bool) ([]*DiscoveryEntry, error)

	Add(key *ServiceKey, params RawParamValues, alias string) error
}

// Services provides access to ServiceRegistry implementation (set in services module)
var Services ServiceRegistry
