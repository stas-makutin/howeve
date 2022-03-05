package defs

import (
	"errors"
	"io"
	"strings"

	"github.com/stas-makutin/howeve/api"
)

// TransportInfo transport definition structure
type TransportInfo struct {
	Name string
	Params
}

// Transport interface - blocking transport operations
type Transport interface {
	ID() api.TransportIdentifier
	Open(entry string, params api.ParamValues) error
	ReadyToRead() <-chan struct{}
	io.ReadWriteCloser
}

// ErrNotOpen error
var ErrNotOpen error = errors.New("the transport entry is not open")

// Transports contains transports definitions (defined in services module)
var Transports map[api.TransportIdentifier]*TransportInfo

// TransportName return name of the transport for provided identifier
func TransportName(t api.TransportIdentifier) string {
	if ti, ok := Transports[t]; ok {
		return ti.Name
	}
	return ""
}

// TransportByName resolves transport name into identifier
func TransportByName(name string) (api.TransportIdentifier, bool) {
	for id, ti := range Transports {
		if strings.EqualFold(name, ti.Name) {
			return id, true
		}
	}
	return 0, false
}
