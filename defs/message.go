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
	Time    time.Time    `json:"time"`
	ID      uuid.UUID    `json:"id"`
	State   MessageState `json:"state"`
	Payload []byte       `json:"payload"`
}

// MessageLog defines message log interface
type MessageLog interface {
	Persist()
	Register(key *ServiceKey, payload []byte, state MessageState) *Message
	UpdateState(id uuid.UUID, state MessageState) (*ServiceKey, *Message)
}

// Messages provides access to MessageLog implementation (set in messages module)
var Messages MessageLog
