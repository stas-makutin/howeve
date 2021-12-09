package messages

import (
	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/defs"
)

// message struct defines entry in the message log
type message struct {
	*defs.ServiceKey
	*defs.Message
}

// messages log container
type messages struct {
	services    map[defs.ServiceKey]int
	entries     []*message
	entriesById map[uuid.UUID]*message
}

func newMessages() *messages {
	m := &messages{}
	m.clear()
	return m
}

func (m *messages) clear() {
	m.services = make(map[defs.ServiceKey]int)
	m.entries = nil
	m.entriesById = make(map[uuid.UUID]*message)
}

// pushes new message, returns true if this is first message for provided service or false otherwise
func (m *messages) push(key *defs.ServiceKey, msg *defs.Message) bool {
	messagesCount := m.services[*key]
	m.services[*key] = messagesCount + 1

	entry := &message{
		ServiceKey: key,
		Message:    msg,
	}

	m.entriesById[msg.UUID] = entry
	m.entries = append(m.entries, entry)

	return messagesCount == 0
}

// pops oldest message, returns its service key, content, and true if there's no more messages from its service (or false otherwise)
func (m *messages) pop() (*defs.ServiceKey, *defs.Message, bool) {
	if len(m.entries) <= 0 {
		return nil, nil, false
	}

	entry := m.entries[0]
	m.entries = m.entries[1:]
	delete(m.entriesById, entry.UUID)

	messagesCount := m.services[*entry.ServiceKey] - 1
	if messagesCount == 0 {
		delete(m.services, *entry.ServiceKey)
		return entry.ServiceKey, entry.Message, true
	}

	m.services[*entry.ServiceKey] = messagesCount
	return entry.ServiceKey, entry.Message, false
}
