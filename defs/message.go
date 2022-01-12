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
type MessageFunc func(index int, message *Message) bool

// MessageLog defines message log interface
type MessageLog interface {
	Persist()
	Register(key *ServiceKey, payload []byte, state MessageState) *Message
	UpdateState(id uuid.UUID, state MessageState) (*ServiceKey, *Message)
	Get(id uuid.UUID) (*ServiceKey, *Message)
	List(find MessageFindFunc, filter MessageFunc) int

	// non thread safe
	FromIndex(index int, exclusive bool) MessageFindFunc
	FromID(id uuid.UUID, exclusive bool) MessageFindFunc
	FromTime(time time.Time, exclusive bool) MessageFindFunc
}

func UntilIndex(index int, exclusive bool, next MessageFunc) MessageFunc {
	if exclusive {
		index -= 1
	}
	return func(n int, message *Message) bool {
		if n <= index {
			return true
		}
		if next == nil {
			return false
		}
		return next(n, message)
	}
}

func UntilID(id uuid.UUID, exclusive bool, next MessageFunc) MessageFunc {
	if exclusive {
		return func(index int, message *Message) bool {
			if message.ID == id {
				return true
			}
			if next == nil {
				return false
			}
			return next(index, message)
		}
	}
	return func(index int, message *Message) bool {
		if next != nil && next(index, message) {
			return true
		}
		if message.ID == id {
			return true
		}
		return false
	}
}

func UntilTime(time time.Time, exclusive bool, next MessageFunc) MessageFunc {
	return func(index int, message *Message) bool {
		if message.Time.After(time) {
			return true
		}
		if exclusive && message.Time.Equal(time) {
			return true
		}
		if next == nil {
			return false
		}
		return next(index, message)
	}
}

func HasStates(states []MessageState, next MessageFunc) MessageFunc {
	if len(states) <= 0 {
		return next
	}
	statesMap := make(map[MessageState]struct{})
	for _, state := range states {
		statesMap[state] = struct{}{}
	}
	return func(index int, message *Message) bool {
		if _, ok := statesMap[message.State]; !ok {
			return true
		}
		if next == nil {
			return false
		}
		return next(index, message)
	}
}

// Messages provides access to MessageLog implementation (set in messages module)
var Messages MessageLog
