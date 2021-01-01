package serial

import (
	serial "github.com/albenik/go-serial/v2"
	"github.com/stas-makutin/howeve/services"
	"github.com/stas-makutin/howeve/services/servicedef"
)

// Transport struct - serial Transport implementation
type Transport struct {
	port *serial.Port
}

// Open func
func (t *Transport) Open(entry string, params servicedef.ParamValues) (err error) {
	t.Close()

	options := []serial.Option{}
	if v, ok := params[services.ParamNameSerialBaudRate]; ok {
		options = append(options, serial.WithBaudrate(v.(int)))
	}
	if v, ok := params[services.ParamNameSerialDataBits]; ok {
		switch v {
		case "5":
			options = append(options, serial.WithDataBits(5))
		case "6":
			options = append(options, serial.WithDataBits(6))
		case "7":
			options = append(options, serial.WithDataBits(7))
		case "8":
			options = append(options, serial.WithDataBits(8))
		}
	}
	if v, ok := params[services.ParamNameSerialParity]; ok {
		switch v {
		case "none":
			options = append(options, serial.WithParity(serial.NoParity))
		case "odd":
			options = append(options, serial.WithParity(serial.OddParity))
		case "even":
			options = append(options, serial.WithParity(serial.EvenParity))
		case "mark":
			options = append(options, serial.WithParity(serial.MarkParity))
		case "space":
			options = append(options, serial.WithParity(serial.SpaceParity))
		}
	}
	if v, ok := params[services.ParamNameSerialStopBits]; ok {
		switch v {
		case "1":
			options = append(options, serial.WithStopBits(serial.OneStopBit))
		case "1.5":
			options = append(options, serial.WithStopBits(serial.OnePointFiveStopBits))
		case "2":
			options = append(options, serial.WithStopBits(serial.TwoStopBits))
		}
	}
	if v, ok := params[services.ParamNameSerialWriteTimeout]; ok {
		options = append(options, serial.WithWriteTimeout(int(v.(uint32))))
	}

	t.port, err = serial.Open(entry, options...)

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
