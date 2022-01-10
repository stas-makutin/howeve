package httpsrv

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/events"
	"github.com/stas-makutin/howeve/events/handlers"
)

type queryType uint16

const (
	queryUnexpected = queryType(iota)
	queryRestart
	queryRestartResult
	queryGetConfig
	queryGetConfigResult
	queryProtocolList
	queryProtocolListResult
	queryTransportList
	queryTransportListResult
	queryProtocolInfo
	queryProtocolInfoResult
	queryProtocolDiscover
	queryProtocolDiscoverResult
	queryProtocolDiscovery
	queryProtocolDiscoveryResult
	queryProtocolDiscoveryStarted
	queryProtocolDiscoveryFinished
	queryAddService
	queryAddServiceResult
	queryRemoveService
	queryRemoveServiceResult
	queryChangeServiceAlias
	queryChangeServiceAliasResult
	queryServiceStatus
	queryServiceStatusResult
	queryListServices
	queryListServicesResult
	querySendToService
	querySendToServiceResult
	queryGetMessage
	queryGetMessageResult
	queryListMessages
	queryListMessagesResult
	queryNewMessage
	queryDropMessage
	queryUpdateMessageState
)

var queryTypeMap = map[string]queryType{
	"restart": queryRestart, "restartResult": queryRestartResult,
	"getCfg": queryGetConfig, "getConfig": queryGetConfig, "getConfigResult": queryGetConfigResult,
	"protocols": queryProtocolList, "protocolsResult": queryProtocolListResult,
	"transports": queryTransportList, "transportsResult": queryTransportListResult,
	"protocolInfo": queryProtocolInfo, "protocolInfoResult": queryProtocolInfoResult,
	"discover": queryProtocolDiscover, "discoverResult": queryProtocolDiscoverResult,
	"discovery": queryProtocolDiscovery, "discoveryResult": queryProtocolDiscoveryResult,
	"addService": queryAddService, "addServiceResult": queryAddServiceResult,
	"removeService": queryRemoveService, "removeServiceResult": queryRemoveServiceResult,
	"changeServiceAlias": queryChangeServiceAlias, "changeServiceAliasResult": queryChangeServiceAliasResult,
	"serviceStatus": queryServiceStatus, "serviceStatusResult": queryServiceStatusResult,
	"listServices": queryListServices, "listServicesResult": queryListServicesResult,
	"sendTo": querySendToService, "sendToResult": querySendToServiceResult,
	"getMessage": queryGetMessage, "getMessageResult": queryGetMessageResult,
	"messagesList": queryListMessages, "messagesListResult": queryListMessagesResult,
}
var queryNameMap map[queryType]string

func init() {
	queryNameMap = make(map[queryType]string)
	for k, v := range queryTypeMap {
		queryNameMap[v] = k
	}
}

// Query struct
type Query struct {
	Type    queryType
	ID      string
	Payload interface{}
}

// UnmarshalJSON func
func (c *Query) UnmarshalJSON(data []byte) error {
	var env struct {
		Type    string          `json:"q"`
		ID      string          `json:"i,omitempty"`
		Payload json.RawMessage `json:"p,omitempty"`
	}
	err := json.Unmarshal(data, &env)
	if err != nil {
		return err
	}
	t, ok := queryTypeMap[env.Type]
	if !ok {
		return fmt.Errorf("unexpected query %v", env.Type)
	}
	c.Type = t
	c.ID = env.ID
	return c.unmarshalPayload(env.Payload)
}

// MarshalJSON func
func (c *Query) MarshalJSON() ([]byte, error) {
	var env struct {
		Type    string      `json:"q"`
		ID      string      `json:"i,omitempty"`
		Payload interface{} `json:"p,omitempty"`
	}
	n, ok := queryNameMap[c.Type]
	if !ok {
		return nil, fmt.Errorf("unknown query %v", c.Type)
	}
	env.Type = n
	env.ID = c.ID
	env.Payload = c.Payload
	return json.Marshal(env)
}

func (c *Query) unmarshalPayload(data []byte) error {
	switch c.Type {
	case queryProtocolInfo:
		if len(data) > 0 {
			var p handlers.ProtocolInfoFilter
			if err := json.Unmarshal(data, &p); err != nil {
				return err
			}
			c.Payload = &p
		}
	case queryProtocolDiscover:
		var p handlers.ProtocolDiscoverInput
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case queryProtocolDiscovery:
		var p handlers.ProtocolDiscoveryInput
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case queryAddService:
		var p handlers.ServiceEntry
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case queryRemoveService:
		var p handlers.ServiceID
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case queryChangeServiceAlias:
		var p handlers.ChangeServiceAliasQuery
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case queryServiceStatus:
		var p handlers.ServiceID
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case queryListServices:
		var p handlers.ListServicesInput
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case querySendToService:
		var p handlers.SendToServiceInput
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case queryGetMessage:
		var p uuid.UUID
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case queryListMessages:
		var p handlers.ListMessagesInput
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	}
	return nil
}

func (c *Query) toEvent() interface{} {
	switch c.Type {
	case queryRestart:
		return &handlers.Restart{RequestHeader: *handlers.NewRequestHeader(c.ID)}
	case queryGetConfig:
		return &handlers.ConfigGet{RequestHeader: *handlers.NewRequestHeader(c.ID)}
	case queryProtocolList:
		return &handlers.ProtocolList{RequestHeader: *handlers.NewRequestHeader(c.ID)}
	case queryTransportList:
		return &handlers.TransportList{RequestHeader: *handlers.NewRequestHeader(c.ID)}
	case queryProtocolInfo:
		var filter *handlers.ProtocolInfoFilter
		if c.Payload != nil {
			filter = c.Payload.(*handlers.ProtocolInfoFilter)
		}
		return &handlers.ProtocolInfo{RequestHeader: *handlers.NewRequestHeader(c.ID), Filter: filter}
	case queryProtocolDiscover:
		return &handlers.ProtocolDiscover{RequestHeader: *handlers.NewRequestHeader(c.ID), ProtocolDiscoverInput: c.Payload.(*handlers.ProtocolDiscoverInput)}
	case queryProtocolDiscovery:
		return &handlers.ProtocolDiscovery{RequestHeader: *handlers.NewRequestHeader(c.ID), ProtocolDiscoveryInput: c.Payload.(*handlers.ProtocolDiscoveryInput)}
	case queryAddService:
		return &handlers.AddService{RequestHeader: *handlers.NewRequestHeader(c.ID), ServiceEntry: c.Payload.(*handlers.ServiceEntry)}
	case queryRemoveService:
		return &handlers.RemoveService{RequestHeader: *handlers.NewRequestHeader(c.ID), ServiceID: c.Payload.(*handlers.ServiceID)}
	case queryChangeServiceAlias:
		return &handlers.ChangeServiceAlias{RequestHeader: *handlers.NewRequestHeader(c.ID), ChangeServiceAliasQuery: c.Payload.(*handlers.ChangeServiceAliasQuery)}
	case queryServiceStatus:
		return &handlers.ServiceStatus{RequestHeader: *handlers.NewRequestHeader(c.ID), ServiceID: c.Payload.(*handlers.ServiceID)}
	case queryListServices:
		return &handlers.ListServices{RequestHeader: *handlers.NewRequestHeader(c.ID), ListServicesInput: c.Payload.(*handlers.ListServicesInput)}
	case querySendToService:
		return &handlers.SendToService{RequestHeader: *handlers.NewRequestHeader(c.ID), SendToServiceInput: c.Payload.(*handlers.SendToServiceInput)}
	case queryGetMessage:
		return &handlers.GetMessage{RequestHeader: *handlers.NewRequestHeader(c.ID), ID: c.Payload.(uuid.UUID)}
	case queryListMessages:
		return &handlers.ListMessages{RequestHeader: *handlers.NewRequestHeader(c.ID), ListMessagesInput: c.Payload.(*handlers.ListMessagesInput)}
	}
	return nil
}

func (c *Query) toTargetedRequest(ctx context.Context, receiverID events.SubscriberID) interface{} {
	event := c.toEvent()
	if te, ok := event.(events.TargetedRequest); ok {
		te.SetReceiver(ctx, receiverID)
		return te
	}
	return event
}

func queryFromEvent(event interface{}) *Query {
	switch e := event.(type) {
	case *handlers.RestartResult:
		return &Query{Type: queryRestartResult, ID: e.TraceID()}
	case *handlers.ConfigGetResult:
		return &Query{Type: queryGetConfigResult, ID: e.TraceID(), Payload: e.Config}
	case *handlers.ProtocolListResult:
		return &Query{Type: queryProtocolListResult, ID: e.TraceID(), Payload: e.Protocols}
	case *handlers.TransportListResult:
		return &Query{Type: queryTransportListResult, ID: e.TraceID(), Payload: e.Transports}
	case *handlers.ProtocolInfoResult:
		return &Query{Type: queryProtocolListResult, ID: e.TraceID(), Payload: e.Protocols}
	case *handlers.ProtocolDiscoverResult:
		return &Query{Type: queryProtocolDiscoverResult, ID: e.TraceID(), Payload: e.ProtocolDiscoverOutput}
	case *handlers.ProtocolDiscoveryResult:
		return &Query{Type: queryProtocolDiscoveryResult, ID: e.TraceID(), Payload: e.DiscoveryResult}
	case *handlers.DiscoveryStarted:
		return &Query{Type: queryProtocolDiscoveryStarted, ID: e.TraceID(), Payload: e.DiscoveryRequest}
	case *handlers.DiscoveryFinished:
		return &Query{Type: queryProtocolDiscoveryFinished, ID: e.TraceID(), Payload: e.DiscoveryResult}
	case *handlers.AddServiceResult:
		return &Query{Type: queryAddServiceResult, ID: e.TraceID(), Payload: e.StatusReply}
	case *handlers.RemoveServiceResult:
		return &Query{Type: queryRemoveServiceResult, ID: e.TraceID(), Payload: e.StatusReply}
	case *handlers.ChangeServiceAliasResult:
		return &Query{Type: queryChangeServiceAliasResult, ID: e.TraceID(), Payload: e.StatusReply}
	case *handlers.ServiceStatusResult:
		return &Query{Type: queryServiceStatusResult, ID: e.TraceID(), Payload: e.StatusReply}
	case *handlers.ListServicesResult:
		return &Query{Type: queryListServicesResult, ID: e.TraceID(), Payload: e.ListServicesOutput}
	case *handlers.SendToServiceResult:
		return &Query{Type: querySendToServiceResult, ID: e.TraceID(), Payload: e.SendToServiceOutput}
	case *handlers.GetMessageResult:
		return &Query{Type: queryGetMessageResult, ID: e.TraceID(), Payload: e.MessageEntry}
	case *handlers.ListMessagesResult:
		return &Query{Type: queryListMessagesResult, ID: e.TraceID(), Payload: e.ListMessagesOutput}
	case *handlers.NewMessage:
		return &Query{Type: queryNewMessage, ID: e.TraceID(), Payload: e.MessageEntry}
	case *handlers.DropMessage:
		return &Query{Type: queryDropMessage, ID: e.TraceID(), Payload: e.MessageEntry}
	case *handlers.UpdateMessageState:
		return &Query{Type: queryUpdateMessageState, ID: e.TraceID(), Payload: e.UpdateMessageStateData}
	}
	return nil
}
