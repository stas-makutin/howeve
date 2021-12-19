package handlers

import (
	"github.com/stas-makutin/howeve/defs"
)

func handleAddService(event *AddService) {
	r := &AddServiceResult{ResponseHeader: event.Associate(), AddServiceReply: &AddServiceReply{Success: false}}
	key, errorInfo := makeServiceKey(event.Protocol, event.Transport, event.Entry)
	if errorInfo == nil {
		if error := defs.Services.Add(key, event.Params, event.Alias); error == nil {
			r.Success = true
		} else {
			switch error {
			case defs.ErrServiceExists:
				errorInfo = NewErrorInfo(ErrorServiceExists,
					defs.ProtocolName(key.Protocol), key.Protocol,
					defs.TransportName(key.Transport), key.Transport, key.Entry,
				)
			case defs.ErrAliasExists:
				errorInfo = NewErrorInfo(ErrorServiceAliasExists, event.Alias)
			default:
				errorInfo = handleParamsErrors(error)
				if errorInfo == nil {
					NewErrorInfo(ErrorServiceInitialize,
						defs.ProtocolName(key.Protocol), key.Protocol,
						defs.TransportName(key.Transport), key.Transport, key.Entry,
						error.Error(),
					)
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

// func handleProtocolDiscovery(event *ProtocolDiscovery) {
// 	r := &ProtocolDiscoveryResult{ResponseHeader: event.Associate(), ProtocolDiscoveryQueryResult: &ProtocolDiscoveryQueryResult{}}
// 	if pat, ei := findProtocolAndTransport(event.Protocol, event.Transport); ei == nil {
// 		if pat.options.DiscoveryFunc != nil {
// 			if params, errorInfo := event.Params.parse(pat.options.DiscoveryParams); errorInfo != nil {
// 				r.ProtocolDiscoveryQueryResult.Error = errorInfo
// 			} else {
// 				go func() {
// 					if serviceEntries, err := pat.options.DiscoveryFunc(event.Context(), params); err == nil {
// 						if len(serviceEntries) > 0 {
// 							r.Services = make([]*ServiceEntryDetails, 0, len(serviceEntries))
// 							for _, serviceEntry := range serviceEntries {
// 								r.Services = append(r.Services, &ServiceEntryDetails{
// 									ServiceEntry: ServiceEntry{
// 										ServiceKey: ServiceKey{
// 											Protocol:  serviceEntry.Key.Protocol,
// 											Transport: serviceEntry.Key.Transport,
// 											Entry:     serviceEntry.Key.Entry,
// 										},
// 										Params: NewParamsValues(serviceEntry.Params),
// 									},
// 									Description: serviceEntry.Description,
// 								})
// 							}
// 						}
// 					} else {
// 						// TODO error reporting
// 					}
// 					Dispatcher.Send(r)
// 				}()
// 				return
// 			}
// 		} else {
// 			r.ProtocolDiscoveryQueryResult.Error = NewErrorInfo(
// 				ErrorNoDiscovery, pat.protocol.Name, event.Protocol, pat.transport.Name, event.Transport,
// 			)
// 		}
// 	} else {
// 		r.ProtocolDiscoveryQueryResult.Error = ei
// 	}

// 	Dispatcher.Send(r)
// }
