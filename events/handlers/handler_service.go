package handlers

import "github.com/stas-makutin/howeve/services"

func handleAddService(event *AddService) {
	r := &AddServiceResult{ResponseHeader: event.Associate(), AddServiceReply: &AddServiceReply{Success: false}}
	entry, errorInfo := makeServiceEntry(event.Protocol, event.Transport, event.Entry, event.Params)
	if errorInfo == nil {
		services.AddService(entry, event.Alias)

	}
	r.Error = errorInfo
	Dispatcher.Send(r)
}

func handleSendToService(event *SendToService) {
	r := &SendToServiceResult{ResponseHeader: event.Associate()}
	Dispatcher.Send(r)
}

func handleRetriveFromService(event *RetriveFromService) {
	r := &RetriveFromServiceResult{ResponseHeader: event.Associate()}
	Dispatcher.Send(r)
}
