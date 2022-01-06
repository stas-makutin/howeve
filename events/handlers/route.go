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
	case *ProtocolDiscover:
		handleProtocolDiscover(e)
	case *ProtocolDiscovery:
		handleProtocolDiscovery(e)
	case *AddService:
		handleAddService(e)
	case *RemoveService:
		handleRemoveService(e)
	case *ChangeServiceAlias:
		handleChangeServiceAlias(e)
	case *ServiceStatus:
		handleServiceStatus(e)
	case *ListServices:
		handleListServices(e)
	case *SendToService:
		handleSendToService(e)
	}
}
