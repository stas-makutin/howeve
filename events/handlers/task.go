package handlers

import (
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/events"
	"github.com/stas-makutin/howeve/tasks"
)

// Dispatcher service event dispatcher
var Dispatcher events.Dispatcher

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
}

func (t *Task) handleEvents(event interface{}) {
	switch e := event.(type) {
	case *Restart:
		handleRestart(e)
	case *ConfigGet:
		handleConfigGet(e, t.cfg)
	case *ProtocolList:
		handleProtocolList(e)
	case *TransportList:
		handleTransportList(e)
	case *ProtocolInfo:
		handleProtocolInfo(e)
	case *ProtocolDiscovery:
		handleProtocolDiscovery(e)
	case *AddService:
		handleAddService(e)
	case *SendToService:
		handleSendToService(e)
	case *RetriveFromService:
		handleRetriveFromService(e)
	}
}
