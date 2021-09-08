package messages

import (
	"sync"
	"time"

	"github.com/stas-makutin/howeve/defs"
)

// message struct defines entry in the message log
type message struct {
	time time.Time
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