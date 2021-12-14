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

var mlogCfg *config.MessageLogConfig
var mlogMaxSize int

var mlog *messages
var mlogSize int

var mlogLock sync.Mutex
var mlogStopWb sync.WaitGroup
var mlogPersistCh chan struct{}
var mlogStopCh chan struct{}

func readConfig(cfg *config.Config, cfgError config.Error) {
	mlogCfg = cfg.MessageLog
	mlogMaxSize = 0
	if mlogCfg != nil {
		mlogMaxSize = int(mlogCfg.MaxSize.Value())
	}
	if mlogMaxSize <= 0 {
		mlogMaxSize = 10 * 1024 * 1024
	} else if mlogMaxSize < 8192 {
		mlogMaxSize = 8192
	} else if mlogMaxSize > 1024*1024*1024 { // 1 GiB
		mlogMaxSize = 1024 * 1024 * 1024
	}
}

func writeConfig(cfg *config.Config) {
	cfg.MessageLog = mlogCfg
}

func openMessageLog() {
	mlog = newMessages()
	mlogSize = minimalLength
	mlogPersistCh = make(chan struct{}, 1)
	mlogStopCh = make(chan struct{})

	mlogStopWb.Add(1)
	go messageLogBackground()
}

func closeMessageLog() {
	mlog = nil
	mlogSize = 0
	close(mlogPersistCh)
	mlogPersistCh = nil
	mlogStopCh = nil
}

func stopMessageLog() {
	close(mlogStopCh)
	mlogStopWb.Wait()
	save()
}

func messageLogBackground() {
	defer mlogStopWb.Done()

	// initial load of message log
	if load() {
		tasks.EndServiceTasks()
		return
	}

	// save cycle
	for {
		select {
		case <-mlogStopCh:
			return
		case <-mlogPersistCh:
		}
		save()
	}
}

func load() bool {
	mlogLock.Lock()
	defer mlogLock.Unlock()
	if mlogCfg != nil && mlogCfg.File != "" {
		var err error
		mlogSize, err = mlog.load(mlogCfg.File, mlogMaxSize)
		if err != nil {
			log.Report(log.SrcMsg, err.Error())
			if mlogCfg.Flags&config.MLFlagIgnoreReadError == 0 {
				return true
			}
		}
	}
	return false
}

func save() {
	mlogLock.Lock()
	defer mlogLock.Unlock()
	if mlogCfg != nil && mlogCfg.File != "" {
		if err := mlog.save(mlogCfg.File, mlogCfg.DirMode.WithDirDefault(), mlogCfg.FileMode.WithFileDefault()); err != nil {
			log.Report(log.SrcMsg, err.Error())
		}
	}
}

// Persist function triggers message log disk saving
func Persist() {
	select {
	case mlogPersistCh <- struct{}{}:
	default:
	}
}

// Register registers new message and add it to the message log
func Register(key *defs.ServiceKey, payload []byte, state defs.MessageState) *defs.Message {
	mlogLock.Lock()
	defer mlogLock.Unlock()

	newSize := mlogSize + messageEntryLength(len(payload))
	if mlog.services[*key] == 0 {
		newSize += serviceEntryLength(len(key.Entry))
	}
	for newSize > mlogMaxSize {
		svc, message, svcLast := mlog.pop()
		if message == nil {
			break
		}
		newSize -= messageEntryLength(len(message.Payload))
		if svcLast {
			newSize -= serviceEntryLength(len(svc.Entry))
		}
		handlers.SendDropMessage(svc, message)
	}
	mlogSize = newSize

	message := &defs.Message{
		Time:    time.Now().UTC(),
		ID:      uuid.New(),
		State:   state,
		Payload: payload,
	}
	mlog.push(key, message)
	handlers.SendNewMessage(key, message)
	return message
}

// UpdateState updates message's state to provided and time to current, by provided message id.
// Returns updated message's service key and content or nil if provided id not found
func UpdateState(id uuid.UUID, state defs.MessageState) (*defs.ServiceKey, *defs.Message) {
	mlogLock.Lock()
	defer mlogLock.Unlock()
	entry := mlog.findByID(id)
	if entry != nil {
		if index, found := mlog.findByTime(entry.Time); found {
			length := len(mlog.entries)
			fentry := mlog.entries[index]
			for {
				if fentry.ID == entry.ID {
					prevState := entry.State
					entry.Time = time.Now().UTC()
					entry.State = state
					copy(mlog.entries[index:], mlog.entries[index+1:])
					mlog.entries[length-1] = entry
					handlers.SendUpdateMessageState(entry.ServiceKey, entry.Message, prevState)
					return entry.ServiceKey, entry.Message
				}
				index++
				fentry = mlog.entries[index]
				if index >= length || !fentry.Time.Equal(entry.Time) {
					break
				}
			}
		}
	}
	return nil, nil
}
