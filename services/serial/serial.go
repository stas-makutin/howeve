package serial

import (
	serial "github.com/albenik/go-serial/v2"
	"github.com/stas-makutin/howeve/services/servicedef"
)

// Transport struct - serial Transport implementation
type Transport struct {
	port *serial.Port
}

// Open func
func (t *Transport) Open(entry string, params servicedef.ParamValues) (err error) {
	t.Close()
	t.port, err = serial.Open(entry, serial.WithBaudrate(115200))
	return
}

// Close func
func (t *Transport) Close() (err error) {
	if t.port != nil {
		err = t.port.Close()
		t.port = nil
	}
	return
}

// Read func
func (t *Transport) Read(p []byte) (int, error) {
	if t.port != nil {
		return t.port.Read(p)
	}
	return 0, servicedef.ErrNotOpen
}

// Write func
func (t *Transport) Write(p []byte) (int, error) {
	if t.port != nil {
		return t.port.Write(p)
	}
	return 0, servicedef.ErrNotOpen
}
