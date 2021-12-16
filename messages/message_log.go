package messages

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/defs"
	"github.com/stas-makutin/howeve/events/handlers"
	"github.com/stas-makutin/howeve/log"
	"github.com/stas-makutin/howeve/tasks"
)

type messageLog struct {
	log  *messages
	size int

	cfg     *config.MessageLogConfig
	maxSize int

	lock      sync.Mutex
	stopWb    sync.WaitGroup
	persistCh chan struct{}
	stopCh    chan struct{}
}

func newMessageLog() *messageLog {
	return &messageLog{}
}

func (ml *messageLog) open() {
	ml.log = newMessages()
	ml.size = minimalLength
	ml.persistCh = make(chan struct{}, 1)
	ml.stopCh = make(chan struct{})

	ml.stopWb.Add(1)
	go ml.persistLoop()
}

func (ml *messageLog) close() {
	ml.log = nil
	ml.size = 0
	close(ml.persistCh)
	ml.persistCh = nil
	ml.stopCh = nil
}

func (ml *messageLog) stop() {
	close(ml.stopCh)
	ml.stopWb.Wait()
	ml.save()
}

func (ml *messageLog) readConfig(cfg *config.Config, cfgError config.Error) {
	ml.cfg = cfg.MessageLog
	ml.maxSize = 0
	if ml.cfg != nil {
		ml.maxSize = int(ml.cfg.MaxSize.Value())
	}
	if ml.maxSize <= 0 {
		ml.maxSize = 10 * 1024 * 1024
	} else if ml.maxSize < 8192 {
		ml.maxSize = 8192
	} else if ml.maxSize > 1024*1024*1024 { // 1 GiB
		ml.maxSize = 1024 * 1024 * 1024
	}
}

func (ml *messageLog) writeConfig(cfg *config.Config) {
	cfg.MessageLog = ml.cfg
}

func (ml *messageLog) persistLoop() {
	defer ml.stopWb.Done()

	// initial load of message log
	if ml.load() {
		tasks.EndServiceTasks()
		return
	}

	// save cycle
	for {
		select {
		case <-ml.stopCh:
			return
		case <-ml.persistCh:
		}
		ml.save()
	}
}

func (ml *messageLog) load() bool {
	ml.lock.Lock()
	defer ml.lock.Unlock()
	if ml.cfg != nil && ml.cfg.File != "" {
		var err error
		ml.size, err = ml.log.load(ml.cfg.File, ml.maxSize)
		if err != nil {
			log.Report(log.SrcMsg, err.Error())
			if ml.cfg.Flags&config.MLFlagIgnoreReadError == 0 {
				return true
			}
		}
	}
	return false
}

func (ml *messageLog) save() {
	ml.lock.Lock()
	defer ml.lock.Unlock()
	if ml.cfg != nil && ml.cfg.File != "" {
		if err := ml.log.save(ml.cfg.File, ml.cfg.DirMode.WithDirDefault(), ml.cfg.FileMode.WithFileDefault()); err != nil {
			log.Report(log.SrcMsg, err.Error())
		}
	}
}

// Persist function triggers message log disk saving
func (ml *messageLog) Persist() {
	select {
	case ml.persistCh <- struct{}{}:
	default:
	}
}

// Register registers new message and add it to the message log
func (ml *messageLog) Register(key *defs.ServiceKey, payload []byte, state defs.MessageState) *defs.Message {
	ml.lock.Lock()
	defer ml.lock.Unlock()

	newSize := ml.size + messageEntryLength(len(payload))
	if ml.log.services[*key] == 0 {
		newSize += serviceEntryLength(len(key.Entry))
	}
	for newSize > ml.maxSize {
		svc, message, svcLast := ml.log.pop()
		if message == nil {
			break
		}
		newSize -= messageEntryLength(len(message.Payload))
		if svcLast {
			newSize -= serviceEntryLength(len(svc.Entry))
		}
		handlers.SendDropMessage(svc, message)
	}
	ml.size = newSize

	message := &defs.Message{
		Time:    time.Now().UTC(),
		ID:      uuid.New(),
		State:   state,
		Payload: payload,
	}
	ml.log.push(key, message)
	handlers.SendNewMessage(key, message)
	return message
}

// UpdateState updates message's state to provided and time to current, by provided message id.
// Returns updated message's service key and content or nil if provided id not found
func (ml *messageLog) UpdateState(id uuid.UUID, state defs.MessageState) (*defs.ServiceKey, *defs.Message) {
	ml.lock.Lock()
	defer ml.lock.Unlock()
	entry := ml.log.findByID(id)
	if entry != nil {
		if index, found := ml.log.findByTime(entry.Time); found {
			length := len(ml.log.entries)
			fentry := ml.log.entries[index]
			for {
				if fentry.ID == entry.ID {
					prevState := entry.State
					entry.Time = time.Now().UTC()
					entry.State = state
					copy(ml.log.entries[index:], ml.log.entries[index+1:])
					ml.log.entries[length-1] = entry
					handlers.SendUpdateMessageState(entry.ServiceKey, entry.Message, prevState)
					return entry.ServiceKey, entry.Message
				}
				index++
				fentry = ml.log.entries[index]
				if index >= length || !fentry.Time.Equal(entry.Time) {
					break
				}
			}
		}
	}
	return nil, nil
}
