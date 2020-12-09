package servicedef

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
