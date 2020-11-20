package services

// ProtocolIdentifier type
type ProtocolIdentifier uint8

// Supported protocols identifiers
const (
	ProtocolZWave = ProtocolIdentifier(iota)
)

// ProtocolInfo protocol definition structure
type ProtocolInfo struct {
	Name       string
	Transports []TransportIdentifier
}

// Protocols contains protocols definitions
var Protocols = map[ProtocolIdentifier]ProtocolInfo{
	ProtocolZWave: {
		Name:       "Z-Wave",
		Transports: []TransportIdentifier{TransportSerial},
	},
}
