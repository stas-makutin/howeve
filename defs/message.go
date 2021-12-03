package defs

import (
	"time"

	"github.com/google/uuid"
)

// MessageState type
type MessageState uint8

// Supported protocols identifiers
const (
	Incoming = MessageState(iota)
	Outgoing
	OutgoingPending
	OutgoingFailed
	OutgoingRejected
)

// Message struct represent the message sent to the service
type Message struct {
	Time    time.Time
	UUID    uuid.UUID
	State   MessageState
	Payload []byte
}
