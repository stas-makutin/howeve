package defs

import (
	"context"
	"errors"
	"strings"

	"github.com/stas-makutin/howeve/api"
)

// ServiceFunc is a method which creates service or returns error
type ServiceFunc func(entry string, params api.ParamValues) (Service, error)

// DiscoveryEntry - discovery entry - information about service instance detected during discovery
type DiscoveryEntry struct {
	api.ServiceKey
	api.ParamValues `json:"params,omitempty"`
	Description     string `json:"description,omitempty"`
}

// DiscoveryFunc is a method which returns discovered service entries or error
type DiscoveryFunc func(ctx context.Context, params api.ParamValues) ([]*DiscoveryEntry, error)

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
	Transports map[api.TransportIdentifier]*ProtocolTransportOptions
}

// Protocols contains protocols definitions (defined in service module)
var Protocols map[api.ProtocolIdentifier]*ProtocolInfo

// ProtocolName return name of the transport for provided identifier
func ProtocolName(p api.ProtocolIdentifier) string {
	if pi, ok := Protocols[p]; ok {
		return pi.Name
	}
	return ""
}

// ProtocolByName resolves protocol name into identifier
func ProtocolByName(name string) (api.ProtocolIdentifier, bool) {
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
	// ErrTransportNotSupported is the error in case if provided transport is not supported
	ErrTransportNotSupported error = errors.New("the transport is not supported")
	// ErrNoProtocolTransport is the error in case if provided transport is not supported for given protocol
	ErrNoProtocolTransport error = errors.New("the transport is not supported for the protocol")
)

// ResolveProtocolAndTransport resolves protocol-transport pair
func ResolveProtocolAndTransport(p api.ProtocolIdentifier, t api.TransportIdentifier) (*ProtocolTransportOptions, *TransportInfo, error) {
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
		return nil, nil, ErrNoProtocolTransport
	}
	return to, ti, nil
}
