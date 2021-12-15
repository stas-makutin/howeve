package defs

import (
	"errors"
	"io"
	"strings"
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
	ID() TransportIdentifier
	Open(entry string, params ParamValues) error
	ReadyToRead() <-chan struct{}
	io.ReadWriteCloser
}

// ErrNotOpen error
var ErrNotOpen error = errors.New("the transport entry is not open")

// IsValid verifies if protocol identifer is valid
func (transport TransportIdentifier) IsValid() bool {
	return transport == TransportSerial
}

// Transports contains transports definitions (defined in services module)
var Transports map[TransportIdentifier]*TransportInfo

// TransportName return name of the transport for provided identifier
func TransportName(t TransportIdentifier) string {
	if ti, ok := Transports[t]; ok {
		return ti.Name
	}
	return ""
}

// TransportByName resolves transport name into identifier
func TransportByName(name string) (TransportIdentifier, bool) {
	for id, ti := range Transports {
		if strings.EqualFold(name, ti.Name) {
			return id, true
		}
	}
	return 0, false
}
