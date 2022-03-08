package httpsrv

import (
	"context"

	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/events"
	"github.com/stas-makutin/howeve/events/handlers"
)

func queryToEvent(c *api.Query) interface{} {
	switch c.Type {
	case api.QueryRestart:
		return &handlers.Restart{RequestHeader: *handlers.NewRequestHeader(c.ID)}
	case api.QueryGetConfig:
		return &handlers.ConfigGet{RequestHeader: *handlers.NewRequestHeader(c.ID)}
	case api.QueryProtocolList:
		return &handlers.ProtocolList{RequestHeader: *handlers.NewRequestHeader(c.ID)}
	case api.QueryTransportList:
		return &handlers.TransportList{RequestHeader: *handlers.NewRequestHeader(c.ID)}
	case api.QueryProtocolInfo:
		var payload *api.ProtocolInfo
		if c.Payload != nil {
			payload = c.Payload.(*api.ProtocolInfo)
		}
		return &handlers.ProtocolInfo{RequestHeader: *handlers.NewRequestHeader(c.ID), ProtocolInfo: payload}
	case api.QueryProtocolDiscover:
		return &handlers.ProtocolDiscover{RequestHeader: *handlers.NewRequestHeader(c.ID), ProtocolDiscover: c.Payload.(*api.ProtocolDiscover)}
	case api.QueryProtocolDiscovery:
		return &handlers.ProtocolDiscovery{RequestHeader: *handlers.NewRequestHeader(c.ID), ProtocolDiscovery: c.Payload.(*api.ProtocolDiscovery)}
	case api.QueryAddService:
		return &handlers.AddService{RequestHeader: *handlers.NewRequestHeader(c.ID), ServiceEntry: c.Payload.(*api.ServiceEntry)}
	case api.QueryRemoveService:
		return &handlers.RemoveService{RequestHeader: *handlers.NewRequestHeader(c.ID), ServiceID: c.Payload.(*api.ServiceID)}
	case api.QueryChangeServiceAlias:
		return &handlers.ChangeServiceAlias{RequestHeader: *handlers.NewRequestHeader(c.ID), ChangeServiceAlias: c.Payload.(*api.ChangeServiceAlias)}
	case api.QueryServiceStatus:
		return &handlers.ServiceStatus{RequestHeader: *handlers.NewRequestHeader(c.ID), ServiceID: c.Payload.(*api.ServiceID)}
	case api.QueryListServices:
		return &handlers.ListServices{RequestHeader: *handlers.NewRequestHeader(c.ID), ListServices: c.Payload.(*api.ListServices)}
	case api.QuerySendToService:
		return &handlers.SendToService{RequestHeader: *handlers.NewRequestHeader(c.ID), SendToService: c.Payload.(*api.SendToService)}
	case api.QueryGetMessage:
		return &handlers.GetMessage{RequestHeader: *handlers.NewRequestHeader(c.ID), ID: c.Payload.(uuid.UUID)}
	case api.QueryListMessages:
		return &handlers.ListMessages{RequestHeader: *handlers.NewRequestHeader(c.ID), ListMessages: c.Payload.(*api.ListMessages)}
	}
	return nil
}

func queryToTargetedRequest(c *api.Query, ctx context.Context, receiverID events.SubscriberID) interface{} {
	event := queryToEvent(c)
	if te, ok := event.(events.TargetedRequest); ok {
		te.SetReceiver(ctx, receiverID)
		return te
	}
	return event
}

func queryFromEvent(event interface{}) *api.Query {
	switch e := event.(type) {
	case *handlers.RestartResult:
		return &api.Query{Type: api.QueryRestartResult, ID: e.TraceID()}
	case *handlers.ConfigGetResult:
		return &api.Query{Type: api.QueryGetConfigResult, ID: e.TraceID(), Payload: e.Config}
	case *handlers.ProtocolListResult:
		return &api.Query{Type: api.QueryProtocolListResult, ID: e.TraceID(), Payload: e.ProtocolListResult}
	case *handlers.TransportListResult:
		return &api.Query{Type: api.QueryTransportListResult, ID: e.TraceID(), Payload: e.TransportListResult}
	case *handlers.ProtocolInfoResult:
		return &api.Query{Type: api.QueryProtocolListResult, ID: e.TraceID(), Payload: e.ProtocolInfoResult}
	case *handlers.ProtocolDiscoverResult:
		return &api.Query{Type: api.QueryProtocolDiscoverResult, ID: e.TraceID(), Payload: e.ProtocolDiscoverResult}
	case *handlers.ProtocolDiscoveryResult:
		return &api.Query{Type: api.QueryProtocolDiscoveryResult, ID: e.TraceID(), Payload: e.ProtocolDiscoveryResult}
	case *handlers.ProtocolDiscoveryStarted:
		return &api.Query{Type: api.QueryProtocolDiscoveryStarted, ID: e.TraceID(), Payload: e.ProtocolDiscoveryStarted}
	case *handlers.ProtocolDiscoveryFinished:
		return &api.Query{Type: api.QueryProtocolDiscoveryFinished, ID: e.TraceID(), Payload: e.ProtocolDiscoveryResult}
	case *handlers.AddServiceResult:
		return &api.Query{Type: api.QueryAddServiceResult, ID: e.TraceID(), Payload: e.StatusReply}
	case *handlers.RemoveServiceResult:
		return &api.Query{Type: api.QueryRemoveServiceResult, ID: e.TraceID(), Payload: e.StatusReply}
	case *handlers.ChangeServiceAliasResult:
		return &api.Query{Type: api.QueryChangeServiceAliasResult, ID: e.TraceID(), Payload: e.StatusReply}
	case *handlers.ServiceStatusResult:
		return &api.Query{Type: api.QueryServiceStatusResult, ID: e.TraceID(), Payload: e.StatusReply}
	case *handlers.ListServicesResult:
		return &api.Query{Type: api.QueryListServicesResult, ID: e.TraceID(), Payload: e.ListServicesResult}
	case *handlers.SendToServiceResult:
		return &api.Query{Type: api.QuerySendToServiceResult, ID: e.TraceID(), Payload: e.SendToServiceResult}
	case *handlers.GetMessageResult:
		return &api.Query{Type: api.QueryGetMessageResult, ID: e.TraceID(), Payload: e.MessageEntry}
	case *handlers.ListMessagesResult:
		return &api.Query{Type: api.QueryListMessagesResult, ID: e.TraceID(), Payload: e.ListMessagesResult}
	case *handlers.NewMessage:
		return &api.Query{Type: api.QueryNewMessage, ID: e.TraceID(), Payload: e.MessageEntry}
	case *handlers.DropMessage:
		return &api.Query{Type: api.QueryDropMessage, ID: e.TraceID(), Payload: e.MessageEntry}
	case *handlers.UpdateMessageState:
		return &api.Query{Type: api.QueryUpdateMessageState, ID: e.TraceID(), Payload: e.UpdateMessageState}
	}
	return nil
}
