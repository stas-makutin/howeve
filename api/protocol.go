package api

// ProtocolIdentifier type
type ProtocolIdentifier uint8

// Supported protocols identifiers
const (
	ProtocolZWave = ProtocolIdentifier(iota + 1)
)

// IsValid verifies if protocol identifer is valid
func (protocol ProtocolIdentifier) IsValid() bool {
	return protocol == ProtocolZWave
}

// TransportIdentifier type
type TransportIdentifier uint8

// Supported transport identifiers
const (
	TransportSerial = TransportIdentifier(iota + 1)
)

// IsValid verifies if protocol identifer is valid
func (transport TransportIdentifier) IsValid() bool {
	return transport == TransportSerial
}
