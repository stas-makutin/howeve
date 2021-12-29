package handlers

import (
	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/defs"
)

func handleAddService(event *AddService) {
	r := &AddServiceResult{ResponseHeader: event.Associate(), AddServiceReply: &AddServiceReply{Success: false}}
	key, errorInfo := makeServiceKey(event.Protocol, event.Transport, event.Entry)
	if errorInfo == nil {
		if err := defs.Services.Add(key, event.Params, event.Alias); err == nil {
			r.Success = true
		} else {
			switch err {
			case defs.ErrServiceExists:
				errorInfo = newErrorInfo(ErrorServiceExists, err, key.Protocol, key.Transport, key.Entry)
			case defs.ErrAliasExists:
				errorInfo = newErrorInfo(ErrorServiceAliasExists, err, event.Alias)
			default:
				errorInfo = handleProtocolErrors(err, key.Protocol, key.Transport)
				if errorInfo == nil {
					errorInfo = handleParamsErrors(err)
				}
				if errorInfo == nil {
					errorInfo = newErrorInfo(ErrorServiceInitialize, err, key.Protocol, key.Transport, key.Entry)
				}
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

func SendDiscoveryStarted(id uuid.UUID, protocol defs.ProtocolIdentifier, transport defs.TransportIdentifier, params defs.RawParamValues) {
	Dispatcher.SendAsync(&DiscoveryStarted{
		Header:    *NewHeader(""),
		ID:        id,
		Protocol:  protocol,
		Transport: transport,
		Params:    params,
	})
}

func SendDiscoveryFinished(id uuid.UUID, entries []*defs.DiscoveryEntry, err error) {
	Dispatcher.SendAsync(&DiscoveryFinished{
		Header: *NewHeader(""),
		DiscoveryResult: &DiscoveryResult{
			ID:      id,
			Entries: entries,
			Error:   newErrorInfo(ErrorDiscoveryFailed, err),
		},
	})
}

func handleProtocolDiscover(event *ProtocolDiscover) {
	r := &ProtocolDiscoverResult{ResponseHeader: event.Associate(), ProtocolDiscoverOutput: &ProtocolDiscoverOutput{}}
	id, err := defs.Services.Discover(event.Protocol, event.Transport, event.Params)
	if id != uuid.Nil {
		r.ID = &id
	}
	if err != nil {
		var errorInfo *ErrorInfo
		switch err {
		case defs.ErrNoDiscovery:
			errorInfo = newErrorInfo(ErrorNoDiscovery, err, event.Protocol, event.Transport)
		case defs.ErrDiscoveryBusy:
			errorInfo = newErrorInfo(ErrorDiscoveryBusy, err)
		default:
			errorInfo = handleProtocolErrors(err, event.Protocol, event.Transport)
			if errorInfo != nil {
				errorInfo = handleParamsErrors(err)
			}
		}
		r.Error = errorInfo
	}
	Dispatcher.Send(r)
}

func handleProtocolDiscovery(event *ProtocolDiscovery) {
	r := &ProtocolDiscoveryResult{ResponseHeader: event.Associate()}
	entries, err := defs.Services.Discovery(event.ID, event.Stop)
	r.Entries = entries
	if err != nil {
		var errorInfo *ErrorInfo
		switch err {
		case defs.ErrNoDiscoveryID:
			errorInfo = newErrorInfo(ErrorNoDiscoveryID, err, event.ID)
		case defs.ErrDiscoveryPending:
			errorInfo = newErrorInfo(ErrorDiscoveryPending, err, event.ID)
		}
		r.Error = errorInfo
	}
	Dispatcher.Send(r)
}
