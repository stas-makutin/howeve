package messages

import (
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/tasks"
)

// Task struct
type Task struct {
	cfg *config.Config
	ml  messages
}

// NewTask func
func NewTask() *Task {
	t := &Task{}
	config.AddReader(t.readConfig)
	config.AddWriter(t.writeConfig)
	return t
}

func (t *Task) readConfig(cfg *config.Config, cfgError config.Error) {
	t.cfg = cfg
}

func (t *Task) writeConfig(cfg *config.Config) {
}

// Open func
func (t *Task) Open(ctx *tasks.ServiceTaskContext) error {
	if t.cfg.MessageLog != nil {

	}
	return nil
}

// Close func
func (t *Task) Close(ctx *tasks.ServiceTaskContext) error {
	return nil
}

// Stop func
func (t *Task) Stop(ctx *tasks.ServiceTaskContext) {
}
