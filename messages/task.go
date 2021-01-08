package messages

import (
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/defs"
	"github.com/stas-makutin/howeve/tasks"
)

// Task struct
type Task struct {
	cfg *defs.Config
}

// NewTask func
func NewTask() *Task {
	t := &Task{}
	config.AddReader(t.readConfig)
	config.AddWriter(t.writeConfig)
	return t
}

func (t *Task) readConfig(cfg *defs.Config, cfgError config.Error) {
	t.cfg = cfg
}

func (t *Task) writeConfig(cfg *defs.Config) {
}

// Open func
func (t *Task) Open(ctx *tasks.ServiceTaskContext) error {
	return nil
}

// Close func
func (t *Task) Close(ctx *tasks.ServiceTaskContext) error {
	return nil
}

// Stop func
func (t *Task) Stop(ctx *tasks.ServiceTaskContext) {
}
