package messages

import (
	"sync"

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
	sync.Mutex
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

func (m *messages) push(key *defs.ServiceKey, msg *defs.Message) {
	m.Lock()
	defer m.Unlock()

	messagesCount := m.services[*key]
	m.services[*key] = messagesCount + 1

	entry := &message{
		ServiceKey: key,
		Message:    msg,
	}

	m.entriesById[msg.UUID] = entry
	m.entries = append(m.entries, entry)
}

func (m *messages) pop() (*defs.ServiceKey, *defs.Message) {
	m.Lock()
	entry := func() *message {
		defer m.Unlock()

		if len(m.entries) <= 0 {
			return nil
		}

		entry := m.entries[0]
		m.entries = m.entries[1:]
		delete(m.entriesById, entry.UUID)
		return entry
	}()

	if entry == nil {
		return nil, nil
	}

	return entry.ServiceKey, entry.Message
}
