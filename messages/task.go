package messages

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/defs"
	"github.com/stas-makutin/howeve/events/handlers"
	"github.com/stas-makutin/howeve/log"
	"github.com/stas-makutin/howeve/tasks"
)

// log operation codes
const (
	opLoad = "L"
	opSave = "S"
)

// messageLog - task implementation
type messageLog struct {
	log  *messages
	size int

	cfg          *api.MessageLogConfig
	maxSize      int
	autoPersists time.Duration

	lock      sync.Mutex
	stopWb    sync.WaitGroup
	persistCh chan struct{}
	stopCh    chan struct{}
}

// NewTask func
func NewTask() *messageLog {
	ml := &messageLog{}
	config.AddReader(ml.readConfig)
	config.AddWriter(ml.writeConfig)
	return ml
}

func (ml *messageLog) Open(ctx *tasks.ServiceTaskContext) error {
	ml.log = newMessages()
	ml.size = minimalLength
	ml.persistCh = make(chan struct{}, 1)
	ml.stopCh = make(chan struct{})

	ml.stopWb.Add(1)
	go ml.persistLoop()

	defs.Messages = ml

	return nil
}

func (ml *messageLog) Close(ctx *tasks.ServiceTaskContext) error {
	defs.Messages = nil

	ml.log = nil
	ml.size = 0
	close(ml.persistCh)
	ml.persistCh = nil
	ml.stopCh = nil
	return nil
}

func (ml *messageLog) Stop(ctx *tasks.ServiceTaskContext) {
	close(ml.stopCh)
	ml.stopWb.Wait()
	ml.save()
}

func (ml *messageLog) readConfig(cfg *api.Config, cfgError config.Error) {
	ml.cfg = cfg.MessageLog

	ml.maxSize = 0
	if ml.cfg != nil {
		ml.maxSize = int(ml.cfg.MaxSize.Value())
		ml.autoPersists = ml.cfg.AutoPesist.Value()
	}
	if ml.maxSize <= 0 {
		ml.maxSize = 10 * 1024 * 1024
	} else if ml.maxSize < 8192 {
		ml.maxSize = 8192
	} else if ml.maxSize > 1024*1024*1024 { // 1 GiB
		ml.maxSize = 1024 * 1024 * 1024
	}

	if ml.autoPersists < 1*time.Second {
		ml.autoPersists = 6 * time.Hour
	}
}

func (ml *messageLog) writeConfig(cfg *api.Config) {
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
		case <-time.After(ml.autoPersists):
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
			log.Report(log.SrcMsg, opLoad, err.Error())
			if ml.cfg.Flags&api.MLFlagIgnoreReadError == 0 {
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
			log.Report(log.SrcMsg, opSave, err.Error())
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
	if index, entry := ml.log.findIndexByID(id); entry != nil {
		prevState := entry.State
		entry.Time = time.Now().UTC()
		entry.State = state
		copy(ml.log.entries[index:], ml.log.entries[index+1:])
		ml.log.entries[len(ml.log.entries)-1] = entry
		handlers.SendUpdateMessageState(entry.ServiceKey, entry.Message, prevState)
		return entry.ServiceKey, entry.Message
	}
	return nil, nil
}

// Get returns single message and associated service key for provided message id
func (ml *messageLog) Get(id uuid.UUID) (*defs.ServiceKey, *defs.Message) {
	ml.lock.Lock()
	defer ml.lock.Unlock()
	entry := ml.log.findByID(id)
	if entry != nil {
		return entry.ServiceKey, entry.Message
	}
	return nil, nil
}

// List messages using filter callback starting from message, found by find callback
func (ml *messageLog) List(find defs.MessageFindFunc, filter defs.MessageFunc) int {
	ml.lock.Lock()
	defer ml.lock.Unlock()

	length := len(ml.log.entries)
	index, ok := find()
	if ok {
		for index < length {
			if filter(index, ml.log.entries[index].ServiceKey, ml.log.entries[index].Message) {
				break
			}
			index++
		}
	}
	return length
}

// FromIndex returns function which search for message with provided index (exclusive false) or next index (exclusive true)
func (ml *messageLog) FromIndex(index int, exclusive bool) defs.MessageFindFunc {
	if index < 0 {
		index = len(ml.log.entries) + index
		if index < 0 {
			index = 0
		}
	}
	if exclusive {
		index += 1
	}
	return func() (int, bool) {
		if index < len(ml.log.entries) {
			return index, true
		}
		return 0, false
	}
}

// FromID returns function which search for message with provided id (exclusive false) or next message (exclusive true)
func (ml *messageLog) FromID(id uuid.UUID, exclusive bool) defs.MessageFindFunc {
	return func() (int, bool) {
		if index, entry := ml.log.findIndexByID(id); entry != nil {
			if exclusive {
				index += 1
				if index >= len(ml.log.entries) {
					return 0, false
				}
			}
			return index, true
		}
		return 0, false
	}
}

// FromTime returns function which search for message with time equal (exclusive false) or after provided
func (ml *messageLog) FromTime(time time.Time, exclusive bool) defs.MessageFindFunc {
	return func() (int, bool) {
		length := len(ml.log.entries)
		index, _ := ml.log.findByTime(time)
		for index < length {
			entry := ml.log.entries[index]
			if !exclusive && entry.Time.Equal(time) {
				return index, true
			}
			if entry.Time.After(time) {
				return index, true
			}
			index++
		}
		return 0, false
	}
}
