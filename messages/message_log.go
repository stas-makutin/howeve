package messages

import (
	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/defs"
	"github.com/stas-makutin/howeve/log"
)

var mlog *messages
var mlogSize int
var mlogCfg *config.MessageLogConfig

func readConfig(cfg *config.Config, cfgError config.Error) {
	mlogCfg = cfg.MessageLog
}

func writeConfig(cfg *config.Config) {
	cfg.MessageLog = mlogCfg
}

func initMessageLog() (err error) {
	mlog = newMessages()
	mlogSize = minimalLength

	if mlogCfg != nil && mlogCfg.File != "" {
		mlogSize, err = mlog.load(mlogCfg.File)
		if err != nil && mlogCfg.Flags&config.MLFlagIgnoreReadError != 0 {
			log.Report(log.SrcMsg, err.Error())
			err = nil
		}
	}
	return
}

func destroyMessageLog() {
	mlog = nil
	mlogSize = 0
}

func Persist() {
	if mlogCfg != nil && mlogCfg.File != "" {
		if err := mlog.save(mlogCfg.File, mlogCfg.DirMode, mlogCfg.FileMode); err != nil {
			log.Report(log.SrcMsg, err.Error())
		}
	}
}

func Register(payload []byte) *defs.Message {
	return nil
}

func SetState(uuid uuid.UUID, state defs.MessageState) bool {
	return false
}
