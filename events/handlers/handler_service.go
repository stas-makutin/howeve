package handlers

import (
	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/defs"
)

func handleAddService(event *AddService) {
	r := &AddServiceResult{ResponseHeader: event.Associate(), StatusReply: &StatusReply{Success: false}}
	errorInfo := validateServiceKey(event.ServiceKey)
	if errorInfo == nil {
		if err := defs.Services.Add(event.ServiceKey, event.Params, event.Alias); err == nil {
			r.Success = true
		} else {
			switch err {
			case defs.ErrServiceExists:
				errorInfo = newErrorInfo(ErrorServiceExists, err, event.Protocol, event.Transport, event.Entry)
			case defs.ErrAliasExists:
				errorInfo = newErrorInfo(ErrorServiceAliasExists, err, event.Alias)
			default:
				errorInfo = handleProtocolErrors(err, event.Protocol, event.Transport)
				if errorInfo == nil {
					errorInfo = handleParamsErrors(err)
				}
				if errorInfo == nil {
					errorInfo = newErrorInfo(ErrorServiceInitialize, err, event.Protocol, event.Transport, event.Entry)
				}
			}
		}
	}
	r.Error = errorInfo
	Dispatcher.Send(r)
}

func handleRemoveService(event *RemoveService) {
	r := &RemoveServiceResult{ResponseHeader: event.Associate(), StatusReply: &StatusReply{Success: false}}
	errorInfo := validateServiceID(event.ServiceKey, event.Alias)
	if errorInfo == nil {
		if err := defs.Services.Remove(event.ServiceKey, event.Alias); err == nil {
			r.Success = true
		} else {
			errorInfo = handleServiceNotExistsError(event.ServiceKey, event.Alias)
		}
	}
	r.Error = errorInfo
	Dispatcher.Send(r)
}

func handleChangeServiceAlias(event *ChangeServiceAlias) {
	r := &ChangeServiceAliasResult{ResponseHeader: event.Associate(), StatusReply: &StatusReply{Success: false}}
	errorInfo := validateServiceID(event.ServiceKey, event.Alias)
	if errorInfo == nil {
		if err := defs.Services.Alias(event.ServiceKey, event.Alias, event.NewAlias); err == nil {
			r.Success = true
		} else {
			errorInfo = handleServiceNotExistsError(event.ServiceKey, event.Alias)
		}
	}
	r.Error = errorInfo
	Dispatcher.Send(r)
}

func handleServiceStatus(event *ServiceStatus) {
	r := &ServiceStatusResult{ResponseHeader: event.Associate(), StatusReply: &StatusReply{Success: false}}
	errorInfo := validateServiceID(event.ServiceKey, event.Alias)
	if errorInfo == nil {
		if status, exists := defs.Services.Status(event.ServiceKey, event.Alias); exists {
			if status == nil || status == defs.ErrStatusGood {
				r.Success = true
			} else {
				errorInfo = newErrorInfo(ErrorServiceStatusBad, status)
			}
		} else {
			errorInfo = handleServiceNotExistsError(event.ServiceKey, event.Alias)
		}
	}
	r.Error = errorInfo
	Dispatcher.Send(r)
}

func handleSendToService(event *SendToService) {
	r := &SendToServiceResult{ResponseHeader: event.Associate()}
	errorInfo := validateServiceID(event.ServiceKey, event.Alias)
	if errorInfo == nil {
		if message, err := defs.Services.Send(event.ServiceKey, event.Alias, event.Payload); err == nil {
			r.Message = message
			r.Success = true
		} else {
			switch err {
			case defs.ErrServiceNotExists:
				errorInfo = handleServiceNotExistsError(event.ServiceKey, event.Alias)
			case defs.ErrBadPayload:
				errorInfo = newErrorInfo(ErrorServiceBadPayload, err)
			case defs.ErrSendBusy:
				errorInfo = newErrorInfo(ErrorServiceSendBusy, err)
			default:
				errorInfo = newErrorInfo(ErrorOtherError, err)
			}
		}
	}
	r.Error = errorInfo
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
