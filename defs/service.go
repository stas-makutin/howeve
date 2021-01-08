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
