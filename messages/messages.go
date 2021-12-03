package messages

import (
	"sync"

	"github.com/stas-makutin/howeve/defs"
)

// message struct defines entry in the message log
type message struct {
	*defs.ServiceKey
	*defs.Message
}

// returns byte length of the message entry
// func (m *message) length() int {
// 	return 8 /* time */ + 2 /* service index */ + len(m.UUID) + len(m.Payload) + 2 /* payload size field */
// }

// messages log container
type messages struct {
	sync.Mutex
	services map[defs.ServiceKey]int
	entries  []*message
}

func newMessages() *messages {
	m := &messages{}
	m.clear()
	return m
}

func (m *messages) clear() {
	m.services = make(map[defs.ServiceKey]int)
	m.entries = nil
}

func (m *messages) push(key *defs.ServiceKey, msg *defs.Message) {
	m.Lock()
	defer m.Unlock()

	messagesCount := m.services[*key]
	m.services[*key] = messagesCount + 1

	m.entries = append(m.entries, &message{
		ServiceKey: key,
		Message:    msg,
	})
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
		return entry
	}()

	if entry == nil {
		return nil, nil
	}

	return entry.ServiceKey, entry.Message
}
