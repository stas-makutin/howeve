package servicedef

import (
	"context"
)

// ProtocolIdentifier type
type ProtocolIdentifier uint8

// Supported protocols identifiers
const (
	ProtocolZWave = ProtocolIdentifier(iota + 1)
)

// ServiceEntryDetails - service entry with details
type ServiceEntryDetails struct {
	ServiceEntry
	Description string
}

// DiscoveryFunc is a method which returns discovered service entries or error
type DiscoveryFunc func(ctx context.Context, params ParamValues) ([]*ServiceEntryDetails, error)

// ProtocolTransportOptions defines transport options specific for the protocol
type ProtocolTransportOptions struct {
	Params          // protocol parameters
	DiscoveryFunc   // could be nil
	DiscoveryParams Params
}

// ProtocolInfo protocol definition structure
type ProtocolInfo struct {
	Name       string
	Transports map[TransportIdentifier]*ProtocolTransportOptions
}
