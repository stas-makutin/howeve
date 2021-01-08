package defs

import "time"

// Message struct represent the message sent to the service
type Message struct {
	UUID    [16]byte
	Payload []byte
}

// MessageEntry struct defines entry in the message log
type MessageEntry struct {
	ServiceKey
	Time time.Time
	Message
}
