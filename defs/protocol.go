package defs

import (
	"context"
	"errors"
	"strings"
)

// ProtocolIdentifier type
type ProtocolIdentifier uint8

// Supported protocols identifiers
const (
	ProtocolZWave = ProtocolIdentifier(iota + 1)
)

// ServiceFunc is a method which creates service or returns error
type ServiceFunc func(entry string, params ParamValues) (Service, error)

// DiscoveryEntry - discovery entry - information about service instance detected during discovery
type DiscoveryEntry struct {
	ServiceKey
	ParamValues
	Description string
}

// DiscoveryFunc is a method which returns discovered service entries or error
type DiscoveryFunc func(ctx context.Context, params ParamValues) ([]*DiscoveryEntry, error)

// ProtocolTransportOptions defines transport options specific for the protocol
type ProtocolTransportOptions struct {
	ServiceFunc     // required
	Params          // protocol parameters
	DiscoveryFunc   // could be nil
	DiscoveryParams Params
}

// ProtocolInfo protocol definition structure
type ProtocolInfo struct {
	Name       string
	Transports map[TransportIdentifier]*ProtocolTransportOptions
}

// IsValid verifies if protocol identifer is valid
func (protocol ProtocolIdentifier) IsValid() bool {
	return protocol == ProtocolZWave
}

// Protocols contains protocols definitions (defined in service module)
var Protocols map[ProtocolIdentifier]*ProtocolInfo

// ProtocolName return name of the transport for provided identifier
func ProtocolName(p ProtocolIdentifier) string {
	if pi, ok := Protocols[p]; ok {
		return pi.Name
	}
	return ""
}

// ProtocolByName resolves protocol name into identifier
func ProtocolByName(name string) (ProtocolIdentifier, bool) {
	for id, pi := range Protocols {
		if strings.EqualFold(name, pi.Name) {
			return id, true
		}
	}
	return 0, false
}

// errors
var (
	// ErrProtocolNotSupported is the error in case if provided protocol is not supported
	ErrProtocolNotSupported error = errors.New("the protocol is not supported")
	// ErrTransportNotSupported is the error in case if provided transport is not supported for given protocol
	ErrTransportNotSupported error = errors.New("the transport is not supported")
)

// ResolveProtocolAndTransport resolves protocol-transport pair
func ResolveProtocolAndTransport(p ProtocolIdentifier, t TransportIdentifier) (*ProtocolTransportOptions, *TransportInfo, error) {
	pi := Protocols[p]
	if pi == nil {
		return nil, nil, ErrProtocolNotSupported
	}
	ti := Transports[t]
	if ti == nil {
		return nil, nil, ErrTransportNotSupported
	}
	to := pi.Transports[t]
	if to == nil {
		return nil, nil, ErrTransportNotSupported
	}
	return to, ti, nil
}
