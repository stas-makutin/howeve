package handlers

import (
	"github.com/stas-makutin/howeve/config"
)

func handleConfigGet(event *ConfigGet, cfg *config.Config) {
	Dispatcher.Send(&ConfigData{Config: *cfg, ResponseHeader: event.Associate()})
}
