package eventh

import (
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/events"
)

// ConfigGet - get config event
type ConfigGet struct {
	events.EventWithReceiver
}

// ConfigData - config data event
type ConfigData struct {
	events.EventWithReceiver
	config.Config
}
