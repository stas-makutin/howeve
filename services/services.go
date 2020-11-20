package services

// ServiceKey struct defines service unique identifier/key
type ServiceKey struct {
	protocol  ProtocolIdentifier
	transport TransportIdentifier
	entry     string
}

func startService(key ServiceKey, params map[string]string) {

}

func stopService(key ServiceKey) {

}
