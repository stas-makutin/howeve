package eventh

import "github.com/stas-makutin/howeve/config"

func handleConfigGet(event *ConfigGet, cfg *config.Config) {
	resp := &ConfigData{Config: *cfg}
	resp.SetReceiver(event.Receiver())
	Dispatcher.Send(resp, resp.Receiver())
}
