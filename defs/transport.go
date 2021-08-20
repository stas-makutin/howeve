package defs

import (
	"errors"
	"io"
)

// TransportIdentifier type
type TransportIdentifier uint8

// Supported transport identifiers
const (
	TransportSerial = TransportIdentifier(iota + 1)
)

// TransportInfo transport definition structure
type TransportInfo struct {
	Name string
	Params
}

// Transport interface - blocking transport operations
type Transport interface {
	Open(entry string, params ParamValues) error
	io.ReadWriteCloser
}

// ErrNotOpen error
var ErrNotOpen error = errors.New("the transport entry is not open")

// IsValid verifies if protocol identifer is valid
func (transport TransportIdentifier) IsValid() bool {
	return transport == TransportSerial
}
