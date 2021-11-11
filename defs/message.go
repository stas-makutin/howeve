package defs

import "github.com/google/uuid"

// Message struct represent the message sent to the service
type Message struct {
	UUID    uuid.UUID
	Payload []byte
}
