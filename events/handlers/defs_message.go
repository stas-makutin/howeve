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

// MessageInfo contains single message information
type MessageInfo struct {
	ID   uuid.UUID `json:"id"`
	Time time.Time `json:"time"`
}

// MessagesInfo - message log information
type MessagesInfo struct {
	Count int          `json:"count"`
	First *MessageInfo `json:"first,omitempty"`
	Last  *MessageInfo `json:"last,omitempty"`
}

// GetMessagesInfo defines get message log information request
type GetMessagesInfo struct {
	RequestHeader
}

// GetMessagesInfoResult defines get message log information result
type GetMessagesInfoResult struct {
	ResponseHeader
	*MessagesInfo
}

// MessagesAfterInput - get messages after particular message inputs
type MessagesAfterInput struct {
	ID uuid.UUID `json:"id"`
	*ListMessagesFilters
}

// MessagesAfter - get messages after particular message request
type MessagesAfter struct {
	RequestHeader
	*MessagesAfterInput
}

// MessagesAfterResult - get messages after particular message result
type MessagesAfterResult struct {
	ResponseHeader
	*ListMessagesOutput
}

// ListMessagesFilters defines filters which could be applied to the list of messages
type ListMessagesFilters struct {
	// TODO
}

// ListMessagesInput defines list messages request inputs
type ListMessagesInput struct {
	From *time.Time `json:"from,omitempty"`
	To   *time.Time `json:"to,omitempty"`
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
	Info     *MessagesInfo   `json:"info"`
}

// ListMessagesResult defines list messages result
type ListMessagesResult struct {
	ResponseHeader
	*ListMessagesOutput
}
