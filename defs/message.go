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
	OutgoingTimedOut
)

// Message struct represent the message sent to the service
type Message struct {
	Time    time.Time    `json:"time"`
	ID      uuid.UUID    `json:"id"`
	State   MessageState `json:"state"`
	Payload []byte       `json:"payload"`
}

// MessageFindFunc is a the callback function used in MessageLog to find first message in the List method
type MessageFindFunc func() (int, bool)

// MessageFunc is a the callback function used in MessageLog methods. Returnning true will stop messages iteration
type MessageFunc func(message *Message) bool

// MessageLog defines message log interface
type MessageLog interface {
	Persist()
	Register(key *ServiceKey, payload []byte, state MessageState) *Message
	UpdateState(id uuid.UUID, state MessageState) (*ServiceKey, *Message)
	Get(id uuid.UUID) (*ServiceKey, *Message)
	List(find MessageFindFunc, filter MessageFunc) int

	// non thread safe
	FindByIndex(index int, exclusive bool) MessageFindFunc
	FindByID(id uuid.UUID, exclusive bool) MessageFindFunc
	FindByTime(time time.Time, exclusive bool) MessageFindFunc
}

// Messages provides access to MessageLog implementation (set in messages module)
var Messages MessageLog
