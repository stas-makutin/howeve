package handlers

import (
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/events"
	"github.com/stas-makutin/howeve/tasks"
)

const MaxNumberOfAsyncSenders = 10

// Dispatcher service event dispatcher
var Dispatcher *events.AsyncDispatcher

// Task struct
type Task struct {
	subscriberID events.SubscriberID
	cfg          *config.Config
}

// NewTask func
func NewTask() *Task {
	t := &Task{}
	config.AddReader(t.readConfig)
	return t
}

func (t *Task) readConfig(cfg *config.Config, cfgError config.Error) {
	t.cfg = cfg
}

// Open func
func (t *Task) Open(ctx *tasks.ServiceTaskContext) error {
	Dispatcher = events.NewAsyncDispatcher(MaxNumberOfAsyncSenders)
	t.subscriberID = Dispatcher.Subscribe(t.handleEvents)
	return nil
}

// Close func
func (t *Task) Close(ctx *tasks.ServiceTaskContext) error {
	Dispatcher.Unsubscribe(t.subscriberID)
	return nil
}

// Stop func
func (t *Task) Stop(ctx *tasks.ServiceTaskContext) {
	Dispatcher.Close()
}
