package api

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

// MessageEntry defines message and corresponding service pair
type MessageEntry struct {
	*ServiceKey
	*Message
}

// UpdateMessageState defines payload of UpdateMessageState notification
type UpdateMessageState struct {
	*MessageEntry
	PrevState MessageState
}

// ListMessages defines list messages request payload
type ListMessages struct {
	FromIndex        *int             `json:"fromIndex,omitempty"`
	FromID           *uuid.UUID       `json:"fromId,omitempty"`
	FromTime         *time.Time       `json:"fromTime,omitempty"`
	FromExclusive    bool             `json:"fromExclusive,omitempty"`
	UntilIndex       *int             `json:"untilIndex,omitempty"`
	UntilID          *uuid.UUID       `json:"untilId,omitempty"`
	UntilTime        *time.Time       `json:"untilTime,omitempty"`
	UntilExclusive   bool             `json:"untilExclusive,omitempty"`
	Count            int              `json:"count"`
	CountAfterFilter bool             `json:"countAfterFilter,omitempty"`
	Services         []*ServiceID     `json:"services,omitempty"`
	States           []MessageState   `json:"states,omitempty"`
	Payloads         [][]PayloadMatch `json:"payloads,omitempty"`
}

// ListMessage contains message info for the list messages response
type ListMessage struct {
	ServiceIndex int `json:"serviceIndex"`
	*Message
}

// ListMessagesResult defines list messages result payload
type ListMessagesResult struct {
	Services []*ServiceID   `json:"services,omitempty"`
	Messages []*ListMessage `json:"messages,omitempty"`
	Count    int            `json:"count"`
}
