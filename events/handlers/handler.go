package handlers

import (
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/tasks"
)

func handleRestart(event *Restart) {
	Dispatcher.Send(&RestartResult{ResponseHeader: event.Associate()})
	go tasks.StopServiceTasks()
}

func handleConfigGet(event *ConfigGet, cfg *config.Config) {
	Dispatcher.Send(&ConfigGetResult{Config: *cfg, ResponseHeader: event.Associate()})
}
