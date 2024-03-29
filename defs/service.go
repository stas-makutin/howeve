package defs

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/api"
)

// ServiceStatus describes the status of the service
type ServiceStatus error

// non-error service status
var ErrStatusGood error = errors.New("the service is functioning normally")

// Service interface, defines minimal set of methods the service needs to support
type Service interface {
	Start()
	Stop()
	Status() ServiceStatus
	Send(payload []byte) (*api.Message, error)
}

// errors
var (
	// ErrServiceExists is the error in case if service already exists
	ErrServiceExists error = errors.New("the service already exists")
	// ErrAliasExists is the error in case if service already exists
	ErrAliasExists error = errors.New("the service alias already exists")
	// ErrServiceNotExists returns if service is not exists (Remove, Status methods)
	ErrServiceNotExists error = errors.New("the service not exists")

	// ErrBadPayload returned by Send method in case if message's payload is not valid has no payload
	ErrBadPayload error = errors.New("the message's payload is not valid")
	// ErrSendBusy returned by Send method in case if service is unable to send message at this time
	ErrSendBusy error = errors.New("the service is too busy and unable to send the message")

	// ErrNoDiscovery returned by Discover method in case if service is not providing discovery function
	ErrNoDiscovery error = errors.New("no discovery service")
	// ErrDiscoveryBusy returned by Discover method in case if there are too many discovery requests are running
	ErrDiscoveryBusy error = errors.New("the discovery service is busy")
	// ErrNoDiscoveryID returned by Discovery method if discovery id is not found
	ErrNoDiscoveryID error = errors.New("the discovery id not found")
	// ErrDiscoveryPending returned by Discovery method if requested discovery query is not completed yet
	ErrDiscoveryPending error = errors.New("the discovery is not completed yet")
)

// ParamNameOpenAttemptsInterval parameter name for the time interval between attempts to open (serial port)
const ParamNameOpenAttemptsInterval = "openAttemptsInterval"

// ParamNameOutgoingMaxTTL parameter name for the maximum time to live of outgoing messages
const ParamNameOutgoingMaxTTL = "outgoingMaxTTL"

// ListFunc is a the callback function used in ServiceRegistry List method. Returnning true will stop services iteration
type ListFunc func(key *api.ServiceKey, alias string, status ServiceStatus, params api.ParamValues) bool

// ResolveIDsInput is the input iteration method for ServiceREgistry ResoveIDs
type ResolveIDsInput func() (key *api.ServiceKey, alias string, stop bool)

// ResolveIDsInput is the output method for ServiceREgistry ResoveIDs
type ResolveIDsOutput func(key *api.ServiceKey, alias string)

// ServiceRegistry defines possible operations with services
type ServiceRegistry interface {
	Discover(protocol api.ProtocolIdentifier, transport api.TransportIdentifier, params api.RawParamValues) (uuid.UUID, error)
	Discovery(id uuid.UUID, stop bool) ([]*api.DiscoveryEntry, error)

	Add(key *api.ServiceKey, params api.RawParamValues, alias string) error
	Alias(key *api.ServiceKey, oldAlias string, newAlias string) error
	Remove(key *api.ServiceKey, alias string) error
	Status(key *api.ServiceKey, alias string) (ServiceStatus, bool)
	List(listFn ListFunc)

	ResolveIDs(out ResolveIDsOutput, in ResolveIDsInput)

	Send(key *api.ServiceKey, alias string, payload []byte) (*api.Message, error)
}

// Services provides access to ServiceRegistry implementation (set in services module)
var Services ServiceRegistry
