package services

import (
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/defs"
	"github.com/stas-makutin/howeve/tasks"
)

// Task struct
type Task struct {
	services *servicesRegistry
}

// NewTask func
func NewTask() *Task {
	t := &Task{
		services: newServicesRegistry(),
	}
	config.AddReader(t.services.readConfig)
	config.AddWriter(t.services.writeConfig)
	return t
}

// Open func
func (t *Task) Open(ctx *tasks.ServiceTaskContext) error {
	t.services.open()
	defs.Services = services
	return nil
}

// Close func
func (t *Task) Close(ctx *tasks.ServiceTaskContext) error {
	defs.Services = nil
	t.services.close()
	return nil
}

// Stop func
func (t *Task) Stop(ctx *tasks.ServiceTaskContext) {
	t.services.stop()
}
