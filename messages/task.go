package messages

import (
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/defs"
	"github.com/stas-makutin/howeve/tasks"
)

// Task struct
type Task struct {
	mlog *messageLog
}

// NewTask func
func NewTask() *Task {
	t := &Task{
		mlog: newMessageLog(),
	}
	config.AddReader(t.mlog.readConfig)
	config.AddWriter(t.mlog.writeConfig)
	return t
}

// Open func
func (t *Task) Open(ctx *tasks.ServiceTaskContext) error {
	t.mlog.open()
	defs.Messages = t.mlog
	return nil
}

// Close func
func (t *Task) Close(ctx *tasks.ServiceTaskContext) error {
	defs.Messages = nil
	t.mlog.close()
	return nil
}

// Stop func
func (t *Task) Stop(ctx *tasks.ServiceTaskContext) {
	t.mlog.stop()
}
