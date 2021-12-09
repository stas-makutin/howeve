package messages

import (
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/tasks"
)

// Task struct
type Task struct {
}

// NewTask func
func NewTask() *Task {
	t := &Task{}
	config.AddReader(readConfig)
	config.AddWriter(writeConfig)
	return t
}

// Open func
func (t *Task) Open(ctx *tasks.ServiceTaskContext) error {
	openMessageLog()
	return nil
}

// Close func
func (t *Task) Close(ctx *tasks.ServiceTaskContext) error {
	closeMessageLog()
	return nil
}

// Stop func
func (t *Task) Stop(ctx *tasks.ServiceTaskContext) {
	stopMessageLog()
}
