package api

// ProtocolListEntry - list of supported protocols
type ProtocolListEntry struct {
	ID   ProtocolIdentifier `json:"id"`
	Name string             `json:"name"`
}

// ProtocolListResult - get list of supported protocols response
type ProtocolListResult struct {
	Protocols []*ProtocolListEntry
}
