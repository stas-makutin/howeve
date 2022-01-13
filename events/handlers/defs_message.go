package handlers

import (
	"time"

	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/defs"
)

// MessageEntry defines message and corresponding service pair
type MessageEntry struct {
	*defs.ServiceKey
	*defs.Message
}

// NewMessage event contains information about new message
type NewMessage struct {
	Header
	*MessageEntry
}

// DropMessage event sent when a message gets removed from message log
type DropMessage struct {
	Header
	*MessageEntry
}

// UpdateMessageStateData contains message state change data
type UpdateMessageStateData struct {
	*MessageEntry
	PrevState defs.MessageState
}

// UpdateMessageState event notifies about message state change
type UpdateMessageState struct {
	Header
	*UpdateMessageStateData
}

// GetMessage - get message request
type GetMessage struct {
	RequestHeader
	ID uuid.UUID `json:"id,omitempty"`
}

// GetMessageResult - get message result
type GetMessageResult struct {
	ResponseHeader
	*MessageEntry
}

// ListMessagesInput defines list messages request inputs
type ListMessagesInput struct {
	FromIndex        *int                  `json:"fromIndex,omitempty"`
	FromID           *uuid.UUID            `json:"fromId,omitempty"`
	FromTime         *time.Time            `json:"fromTime,omitempty"`
	FromExclusive    bool                  `json:"fromExclusive,omitempty"`
	UntilIndex       *int                  `json:"untilIndex,omitempty"`
	UntilID          *uuid.UUID            `json:"untilId,omitempty"`
	UntilTime        *time.Time            `json:"untilTime,omitempty"`
	UntilExclusive   bool                  `json:"untilExclusive,omitempty"`
	Count            int                   `json:"count"`
	CountAfterFilter bool                  `json:"countAfterFilter,omitempty"`
	Services         []ServiceID           `json:"services,omitempty"`
	States           []defs.MessageState   `json:"states,omitempty"`
	Payloads         [][]defs.PayloadMatch `json:"payloads,omitempty"`
}

// ListMessages defines list messages request
type ListMessages struct {
	RequestHeader
	*ListMessagesInput
}

// ListMessagesOutput defines list messages outputs
type ListMessagesOutput struct {
	Messages []*MessageEntry `json:"messages,omitempty"` // TODO - split message and service
	Count    int             `json:"count"`
}

// ListMessagesResult defines list messages result
type ListMessagesResult struct {
	ResponseHeader
	*ListMessagesOutput
}
