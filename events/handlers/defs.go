package handlers

import (
	"github.com/stas-makutin/howeve/config"
)

// Restart - restart the service
type Restart struct {
	RequestHeader
}

// RestartResult - restart the service result
type RestartResult struct {
	ResponseHeader
}

// ConfigGet - get config event
type ConfigGet struct {
	RequestHeader
}

// ConfigGetResult - config data event
type ConfigGetResult struct {
	ResponseHeader
	config.Config
}
