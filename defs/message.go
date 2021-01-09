package defs

// Message struct represent the message sent to the service
type Message struct {
	UUID    [16]byte
	Payload []byte
}
