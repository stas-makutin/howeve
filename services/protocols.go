package services

import "context"

// ProtocolIdentifier type
type ProtocolIdentifier uint8

// Supported protocols identifiers
const (
	ProtocolZWave = ProtocolIdentifier(iota)
)

// ServiceKey struct defines service unique identifier/key
type ServiceKey struct {
	protocol  ProtocolIdentifier
	transport TransportIdentifier
	entry     string
}

// ServiceEntry defines service entry - i.e. entry point of service execution
type ServiceEntry struct {
	Key    ServiceKey
	Params ParamValues
}

// DiscoveryFunc is a method which returns discovered service entries or error
type DiscoveryFunc func(ctx context.Context, params ParamValues) ([]ServiceEntry, error)

// ProtocolTransportOptions defines transport options specific for the protocol
type ProtocolTransportOptions struct {
	Params          // protocol parameters
	DiscoveryFunc   // could be nil
	DiscoveryParams Params
}

// ProtocolInfo protocol definition structure
type ProtocolInfo struct {
	Name       string
	Transports map[TransportIdentifier]ProtocolTransportOptions
}

// Protocols contains protocols definitions
var Protocols = map[ProtocolIdentifier]ProtocolInfo{
	ProtocolZWave: {
		Name: "Z-Wave",
		Transports: map[TransportIdentifier]ProtocolTransportOptions{
			TransportSerial: {
				DiscoveryFunc: nil,
			},
		},
	},
}
