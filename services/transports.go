package services

// TransportIdentifier type
type TransportIdentifier uint8

// Supported transport identifiers
const (
	TransportSerial = TransportIdentifier(iota)
)

// Transports contains transports definitions
var Transports = map[TransportIdentifier]string{
	TransportSerial: "Serial",
}
