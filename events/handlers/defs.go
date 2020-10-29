package handlers

import (
	"github.com/stas-makutin/howeve/config"
)

// ConfigGet - get config event
type ConfigGet struct {
	RequestHeader
}

// ConfigData - config data event
type ConfigData struct {
	ResponseHeader
	config.Config
}
