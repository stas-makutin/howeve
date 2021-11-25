package defs

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
