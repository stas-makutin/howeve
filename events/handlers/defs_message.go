package handlers

import (
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

// ListMessagesFilters defines filters which could be applied to the list of messages
type ListMessagesFilters struct {
	// TODO
}

// ListMessagesInput defines list messages request inputs
type ListMessagesInput struct {
	*ListMessagesFilters
}

// ListMessages defines list messages request
type ListMessages struct {
	RequestHeader
	*ListMessagesInput
}

// ListMessagesOutput defines list messages outputs
type ListMessagesOutput struct {
	Messages []*MessageEntry `json:"messages,omitempty"`
	Count    int             `json:"count"`
}

// ListMessagesResult defines list messages result
type ListMessagesResult struct {
	ResponseHeader
	*ListMessagesOutput
}
