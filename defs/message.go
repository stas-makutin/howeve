package defs

import (
	"bytes"
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
type MessageFunc func(index int, key *ServiceKey, message *Message) bool

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

// UntilIndex limits messages iteration by provided index, could be inclusive or exclusive
func UntilIndex(index int, exclusive bool, next MessageFunc) MessageFunc {
	if exclusive {
		index -= 1
	}
	return func(n int, key *ServiceKey, message *Message) bool {
		if n <= index {
			return true
		}
		if next == nil {
			return false
		}
		return next(n, key, message)
	}
}

// UntilID limits messages iteration by provided id, could be inclusive or exclusive
func UntilID(id uuid.UUID, exclusive bool, next MessageFunc) MessageFunc {
	if exclusive {
		return func(index int, key *ServiceKey, message *Message) bool {
			if message.ID == id {
				return true
			}
			if next == nil {
				return false
			}
			return next(index, key, message)
		}
	}
	return func(index int, key *ServiceKey, message *Message) bool {
		if next != nil && next(index, key, message) {
			return true
		}
		if message.ID == id {
			return true
		}
		return false
	}
}

// UntilTime limits messages iteration by provided time, could be inclusive or exclusive
func UntilTime(time time.Time, exclusive bool, next MessageFunc) MessageFunc {
	return func(index int, key *ServiceKey, message *Message) bool {
		if message.Time.After(time) {
			return true
		}
		if exclusive && message.Time.Equal(time) {
			return true
		}
		if next == nil {
			return false
		}
		return next(index, key, message)
	}
}

// UntilCounter limits messages iteration by provided count
func UntilCounter(count int, next MessageFunc) MessageFunc {
	i := 0
	return func(index int, key *ServiceKey, message *Message) bool {
		if i >= count {
			return true
		}
		i += 1
		if next == nil {
			return false
		}
		return next(index, key, message)
	}
}

// WithPayload filters messages based on their state
func WithStates(states []MessageState, next MessageFunc) MessageFunc {
	if len(states) <= 0 {
		return next
	}
	statesMap := make(map[MessageState]struct{})
	for _, state := range states {
		statesMap[state] = struct{}{}
	}
	return func(index int, key *ServiceKey, message *Message) bool {
		if _, ok := statesMap[message.State]; !ok {
			if next != nil {
				return next(index, key, message)
			}
		}
		return false
	}
}

// WithPayload filters messages based on their services
func WithServices(services []*ServiceKey, next MessageFunc) MessageFunc {
	if len(services) <= 0 {
		return next
	}
	return func(index int, key *ServiceKey, message *Message) bool {
		for _, service := range services {
			if *service == *key {
				if next != nil {
					return next(index, key, message)
				}
			}
		}
		return false
	}
}

// PayloadMatch defines byte sequence to match with payload
// At == nil 	- payload must contain the content
// At >= 0 		- payload must include the content at provided position "At"
// At < 0 		- payload must include the content at provided position "len(payload) - len(content) + At + 1"
type PayloadMatch struct {
	Content []byte `json:"content,omitempty"`
	At      *int   `json:"at,omitempty"`
}

// WithPayload filters messages based on their payload
func WithPayload(matches [][]PayloadMatch, next MessageFunc) MessageFunc {
	if len(matches) <= 0 {
		return next
	}
	return func(index int, key *ServiceKey, message *Message) bool {
		for _, match := range matches {
			if len(match) > 0 {
				matches := true
				for _, seq := range match {
					if seq.At == nil {
						if !bytes.Contains(message.Payload, seq.Content) {
							matches = false
							break
						}
					} else {
						pl := len(message.Payload)
						sl := len(seq.Content)
						index := *seq.At
						if index < 0 {
							index = pl - sl + index + 1
							if index < 0 {
								matches = false
								break
							}
						} else if index+sl >= pl {
							matches = false
							break
						}
						if !bytes.Equal(message.Payload[index:index+sl], seq.Content) {
							matches = false
							break
						}
					}
				}
				if matches {
					if next != nil {
						return next(index, key, message)
					}
				}
			}
		}
		return false
	}
}

// Messages provides access to MessageLog implementation (set in messages module)
var Messages MessageLog
