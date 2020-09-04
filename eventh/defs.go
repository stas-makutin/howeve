package eventh

import (
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/events"
)

// ConfigGet - get config event
type ConfigGet struct {
	events.RequestTarget
}

// ConfigData - config data event
type ConfigData struct {
	events.ResponseTarget
	config.Config
}
