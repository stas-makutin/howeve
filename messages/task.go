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
func (t *Task) Open(ctx *tasks.ServiceTaskContext) (err error) {
	err = initMessageLog()
	return
}

// Close func
func (t *Task) Close(ctx *tasks.ServiceTaskContext) error {
	destroyMessageLog()
	return nil
}

// Stop func
func (t *Task) Stop(ctx *tasks.ServiceTaskContext) {
	Persist()
}
