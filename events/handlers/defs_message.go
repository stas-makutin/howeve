package handlers

import "github.com/stas-makutin/howeve/defs"

// NewMessage event contains information about new message
type NewMessage struct {
	Header
	*defs.ServiceKey
	*defs.Message
}

// DropMessage event sent when a message gets removed from message log
type DropMessage struct {
	Header
	*defs.ServiceKey
	*defs.Message
}

// UpdateMessageState event notifies about message state change
type UpdateMessageState struct {
	Header
	PrevState defs.MessageState
	*defs.ServiceKey
	*defs.Message
}
