package defs

import "errors"

// ServiceKey struct defines service unique identifier/key
type ServiceKey struct {
	Protocol  ProtocolIdentifier
	Transport TransportIdentifier
	Entry     string
}

// ServiceEntry defines service entry - i.e. entry point of service execution
type ServiceEntry struct {
	Key    ServiceKey
	Params ParamValues
}

// ServiceStatus describes the status of the service
type ServiceStatus struct {
}

// Service interface, defines minimal set of methods the service needs to support
type Service interface {
	Start()
	Stop()
	Status() ServiceStatus
	Send(message Message) error
}

// errors
var (
	// ErrServiceExists is the error in case if service already exists
	ErrServiceExists error = errors.New("the service already exists")
	// ErrAliasExists is the error in case if service already exists
	ErrAliasExists error = errors.New("the service alias already exists")
	// ErrBadPayload returned by Send method in case if message's payload is not valid has no payload
	ErrBadPayload error = errors.New("the message's payload is not valid")
	// ErrSendBusy returned by Send method in case if service is unable to send message at this time
	ErrSendBusy error = errors.New("the service is too busy and unable to send the message")
)

// ParamNameOpenAttemptsInterval parameter name for the time interval between attempts to open serial port
const ParamNameOpenAttemptsInterval = "openAttemptsInterval"
