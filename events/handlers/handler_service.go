package handlers

import (
	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/defs"
)

func handleAddService(event *AddService) {
	r := &AddServiceResult{ResponseHeader: event.Associate(), StatusReply: &api.StatusReply{Success: false}}
	errorInfo := validateServiceKey(event.ServiceKey)
	if errorInfo == nil {
		if err := defs.Services.Add(event.ServiceKey, event.Params, event.Alias); err == nil {
			r.Success = true
		} else {
			switch err {
			case defs.ErrServiceExists:
				errorInfo = newErrorInfo(api.ErrorServiceExists, err, event.Protocol, event.Transport, event.Entry)
			case defs.ErrAliasExists:
				errorInfo = newErrorInfo(api.ErrorServiceAliasExists, err, event.Alias)
			default:
				errorInfo = handleProtocolErrors(err, event.Protocol, event.Transport)
				if errorInfo == nil {
					errorInfo = handleParamsErrors(err)
				}
				if errorInfo == nil {
					errorInfo = newErrorInfo(api.ErrorServiceInitialize, err, event.Protocol, event.Transport, event.Entry)
				}
			}
		}
	}
	r.Error = errorInfo
	Dispatcher.Send(r)
}

func handleRemoveService(event *RemoveService) {
	r := &RemoveServiceResult{ResponseHeader: event.Associate(), StatusReply: &api.StatusReply{Success: false}}
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
	r := &ChangeServiceAliasResult{ResponseHeader: event.Associate(), StatusReply: &api.StatusReply{Success: false}}
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
	r := &ServiceStatusResult{ResponseHeader: event.Associate(), StatusReply: &api.StatusReply{Success: false}}
	errorInfo := validateServiceID(event.ServiceKey, event.Alias)
	if errorInfo == nil {
		if status, exists := defs.Services.Status(event.ServiceKey, event.Alias); exists {
			if status == nil || status == defs.ErrStatusGood {
				r.Success = true
			} else {
				errorInfo = newErrorInfo(api.ErrorServiceStatusBad, status)
			}
		} else {
			errorInfo = handleServiceNotExistsError(event.ServiceKey, event.Alias)
		}
	}
	r.Error = errorInfo
	Dispatcher.Send(r)
}

func handleListServices(event *ListServices) {
	r := &ListServicesResult{ResponseHeader: event.Associate(), ListServicesResult: &api.ListServicesResult{}}
	defs.Services.List(func(key *api.ServiceKey, alias string, status defs.ServiceStatus, params api.ParamValues) bool {
		found := 0b1111
		if len(event.Protocols) > 0 {
			mask := 0b0001
			found &= ^mask
			for _, protocol := range event.Protocols {
				if key.Protocol == protocol {
					found |= mask
					break
				}
			}
		}
		if len(event.Transports) > 0 {
			mask := 0b0010
			found &= ^mask
			for _, transport := range event.Transports {
				if key.Transport == transport {
					found |= mask
					break
				}
			}
		}
		if len(event.Entries) > 0 {
			mask := 0b0100
			found &= ^mask
			for _, entry := range event.Entries {
				if key.Entry == entry {
					found |= mask
					break
				}
			}
		}
		if len(event.Aliases) > 0 {
			mask := 0b1000
			found &= ^mask
			for _, aliasFilter := range event.Aliases {
				if alias == aliasFilter {
					found |= mask
					break
				}
			}
		}
		if found == 0b1111 {
			statusReply := &api.StatusReply{Success: false}
			if status == nil || status == defs.ErrStatusGood {
				statusReply.Success = true
			} else {
				statusReply.Error = newErrorInfo(api.ErrorServiceStatusBad, status)
			}
			r.Services = append(r.Services, api.ListServicesEntry{
				ServiceEntry: &api.ServiceEntry{ServiceKey: key, Alias: alias, Params: params.Raw()},
				StatusReply:  statusReply,
			})
		}
		return false
	})
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
				errorInfo = newErrorInfo(api.ErrorServiceBadPayload, err)
			case defs.ErrSendBusy:
				errorInfo = newErrorInfo(api.ErrorServiceSendBusy, err)
			default:
				errorInfo = newErrorInfo(api.ErrorOtherError, err)
			}
		}
	}
	r.Error = errorInfo
	Dispatcher.Send(r)
}

func SendDiscoveryStarted(id uuid.UUID, protocol api.ProtocolIdentifier, transport api.TransportIdentifier, params api.RawParamValues) {
	Dispatcher.SendAsync(&ProtocolDiscoveryStarted{
		Header: *NewHeader(""),
		ProtocolDiscoveryStarted: &api.ProtocolDiscoveryStarted{
			ID: id,
			ProtocolDiscover: api.ProtocolDiscover{
				Protocol:  protocol,
				Transport: transport,
				Params:    params,
			},
		},
	})
}

func SendDiscoveryFinished(id uuid.UUID, entries []*api.DiscoveryEntry, err error) {
	Dispatcher.SendAsync(&ProtocolDiscoveryFinished{
		Header: *NewHeader(""),
		ProtocolDiscoveryResult: &api.ProtocolDiscoveryResult{
			ID:      id,
			Entries: entries,
			Error:   newErrorInfo(api.ErrorDiscoveryFailed, err),
		},
	})
}

func handleProtocolDiscover(event *ProtocolDiscover) {
	r := &ProtocolDiscoverResult{ResponseHeader: event.Associate(), ProtocolDiscoverResult: &api.ProtocolDiscoverResult{}}
	id, err := defs.Services.Discover(event.Protocol, event.Transport, event.Params)
	if id != uuid.Nil {
		r.ID = &id
	}
	if err != nil {
		var errorInfo *api.ErrorInfo
		switch err {
		case defs.ErrNoDiscovery:
			errorInfo = newErrorInfo(api.ErrorNoDiscovery, err, event.Protocol, event.Transport)
		case defs.ErrDiscoveryBusy:
			errorInfo = newErrorInfo(api.ErrorDiscoveryBusy, err)
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
		var errorInfo *api.ErrorInfo
		switch err {
		case defs.ErrNoDiscoveryID:
			errorInfo = newErrorInfo(api.ErrorNoDiscoveryID, err, event.ID)
		case defs.ErrDiscoveryPending:
			errorInfo = newErrorInfo(api.ErrorDiscoveryPending, err, event.ID)
		default:
			errorInfo = newErrorInfo(api.ErrorDiscoveryFailed, err)
		}
		r.Error = errorInfo
	}
	Dispatcher.Send(r)
}
