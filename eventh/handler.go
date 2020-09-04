package eventh

import (
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/events"
)

func handleConfigGet(event *ConfigGet, cfg *config.Config) {
	Dispatcher.Send(&ConfigData{Config: *cfg, ResponseTarget: events.ResponseTarget{event.ReceiverID}}, event.ReceiverID)
}
