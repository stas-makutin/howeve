package messages

import (
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/api"
)

// message struct defines entry in the message log
type message struct {
	*api.ServiceKey
	*api.Message
}

// messages log container
type messages struct {
	services    map[api.ServiceKey]int
	entries     []*message
	entriesById map[uuid.UUID]*message
}

func newMessages() *messages {
	m := &messages{}
	m.clear()
	return m
}

func (m *messages) clear() {
	m.services = make(map[api.ServiceKey]int)
	m.entries = nil
	m.entriesById = make(map[uuid.UUID]*message)
}

// pushes new message, returns true if this is first message for provided service or false otherwise
func (m *messages) push(key *api.ServiceKey, msg *api.Message) bool {
	messagesCount := m.services[*key]
	m.services[*key] = messagesCount + 1

	entry := &message{
		ServiceKey: key,
		Message:    msg,
	}

	m.entriesById[msg.ID] = entry
	m.entries = append(m.entries, entry)

	return messagesCount == 0
}

// pops oldest message, returns its service key, content, and true if there's no more messages from its service (or false otherwise)
func (m *messages) pop() (*api.ServiceKey, *api.Message, bool) {
	if len(m.entries) <= 0 {
		return nil, nil, false
	}

	entry := m.entries[0]
	m.entries = m.entries[1:]
	delete(m.entriesById, entry.ID)

	messagesCount := m.services[*entry.ServiceKey] - 1
	if messagesCount == 0 {
		delete(m.services, *entry.ServiceKey)
		return entry.ServiceKey, entry.Message, true
	}

	m.services[*entry.ServiceKey] = messagesCount
	return entry.ServiceKey, entry.Message, false
}

func (m *messages) findByID(id uuid.UUID) *message {
	return m.entriesById[id]
}

func (m *messages) findByTime(t time.Time) (int, bool) {
	length := len(m.entries)
	index := sort.Search(length, func(i int) bool {
		et := m.entries[i].Time
		return et.Equal(t) || et.After(t)
	})
	return index, index < length && m.entries[index].Time.Equal(t)
}

func (m *messages) findIndexByID(id uuid.UUID) (int, *message) {
	entry := m.entriesById[id]
	if entry != nil {
		if index, found := m.findByTime(entry.Time); found {
			length := len(m.entries)
			fentry := m.entries[index]
			for {
				if fentry.ID == entry.ID {
					return index, entry
				}
				index++
				if index >= length {
					break
				}
				fentry = m.entries[index]
				if !fentry.Time.Equal(entry.Time) {
					break
				}
			}
		}
	}
	return 0, nil
}
