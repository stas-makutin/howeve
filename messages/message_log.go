package messages

import (
	"sync"

	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/defs"
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

	mlogMaxSize = int(mlogCfg.MaxSize.Value())
	if mlogMaxSize <= 0 {
		mlogMaxSize = 10 * 1024 * 1024
	} else if mlogMaxSize < 8192 {
		mlogMaxSize = 8192
	} else if mlogMaxSize < 1024*1024*1024 { // 1 GiB
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

func Persist() {
	select {
	case mlogPersistCh <- struct{}{}:
	default:
	}
}

func Register(payload []byte) *defs.Message {
	return nil
}

func SetState(uuid uuid.UUID, state defs.MessageState) bool {
	return false
}
