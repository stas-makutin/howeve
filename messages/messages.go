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
func (m *message) length() int {
	return 8 /* time */ + 2 /* service index */ + len(m.UUID) + len(m.Payload) + 2 /* payload size field */
}

// messages log container
type messages struct {
	sync.Mutex
	services map[defs.ServiceKey]int
	entries  []*message
}

func (m *messages) load(file string) error {
	return nil
}

func (m *messages) save(file string) error {
	return nil
}
