package handlers

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
	}
}
