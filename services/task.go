package services

import (
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/defs"
	"github.com/stas-makutin/howeve/tasks"
)

// Task struct
type Task struct {
	cfg *config.Config
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
	// t.cfg.Services = TODO build configuration
}

// Open func
func (t *Task) Open(ctx *tasks.ServiceTaskContext) error {
	services = newServicesRegistry()
	defs.Services = services
	for _, scfg := range t.cfg.Services {
		addServiceFromConfig(scfg)
	}
	return nil
}

// Close func
func (t *Task) Close(ctx *tasks.ServiceTaskContext) error {
	services = nil
	return nil
}

// Stop func
func (t *Task) Stop(ctx *tasks.ServiceTaskContext) {
	services.Stop()
}
