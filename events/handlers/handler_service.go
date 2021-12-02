package handlers

import (
	"github.com/stas-makutin/howeve/defs"
	"github.com/stas-makutin/howeve/services"
)

func handleAddService(event *AddService) {
	r := &AddServiceResult{ResponseHeader: event.Associate(), AddServiceReply: &AddServiceReply{Success: false}}
	entry, errorInfo := makeServiceEntry(event.Protocol, event.Transport, event.Entry, event.Params)
	if errorInfo == nil {
		if error := services.AddService(entry, event.Alias); error == nil {
			r.Success = true
		} else {
			switch error {
			case defs.ErrServiceExists:
				errorInfo = NewErrorInfo(ErrorServiceExists,
					services.ProtocolName(entry.Key.Protocol), entry.Key.Protocol,
					services.TransportName(entry.Key.Transport), entry.Key.Transport, entry.Key.Entry,
				)
			case defs.ErrAliasExists:
				errorInfo = NewErrorInfo(ErrorServiceAliasExists, event.Alias)
			default:
				errorInfo = NewErrorInfo(ErrorServiceInitialize,
					services.ProtocolName(entry.Key.Protocol), entry.Key.Protocol,
					services.TransportName(entry.Key.Transport), entry.Key.Transport, entry.Key.Entry,
					error.Error(),
				)
			}
		}
	}
	r.Error = errorInfo
	Dispatcher.Send(r)
}

func handleSendToService(event *SendToService) {
	r := &SendToServiceResult{ResponseHeader: event.Associate()}
	Dispatcher.Send(r)
}
