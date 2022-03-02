package api

import (
	"encoding/json"
	"fmt"
)

// QueryType type
type QueryType uint16

// All supported queries
const (
	QueryUnexpected = QueryType(iota)
	QueryRestart
	QueryRestartResult
	QueryGetConfig
	QueryGetConfigResult
	QueryProtocolList
	QueryProtocolListResult
	QueryTransportList
	QueryTransportListResult
	QueryProtocolInfo
	QueryProtocolInfoResult
	QueryProtocolDiscover
	QueryProtocolDiscoverResult
	QueryProtocolDiscovery
	QueryProtocolDiscoveryResult
	QueryProtocolDiscoveryStarted
	QueryProtocolDiscoveryFinished
	QueryAddService
	QueryAddServiceResult
	QueryRemoveService
	QueryRemoveServiceResult
	QueryChangeServiceAlias
	QueryChangeServiceAliasResult
	QueryServiceStatus
	QueryServiceStatusResult
	QueryListServices
	QueryListServicesResult
	QuerySendToService
	QuerySendToServiceResult
	QueryGetMessage
	QueryGetMessageResult
	QueryListMessages
	QueryListMessagesResult
	QueryNewMessage
	QueryDropMessage
	QueryUpdateMessageState
	QueryEventSubscribe
	QueryEventSubscribeResult
)

var queryTypeMap = map[string]QueryType{
	"restart": QueryRestart, "restartResult": QueryRestartResult,
	"getCfg": QueryGetConfig, "getConfig": QueryGetConfig, "getConfigResult": QueryGetConfigResult,
	"protocols": QueryProtocolList, "protocolsResult": QueryProtocolListResult,
	"transports": QueryTransportList, "transportsResult": QueryTransportListResult,
	"protocolInfo": QueryProtocolInfo, "protocolInfoResult": QueryProtocolInfoResult,
	"discover": QueryProtocolDiscover, "discoverResult": QueryProtocolDiscoverResult,
	"discovery": QueryProtocolDiscovery, "discoveryResult": QueryProtocolDiscoveryResult,
	"addService": QueryAddService, "addServiceResult": QueryAddServiceResult,
	"removeService": QueryRemoveService, "removeServiceResult": QueryRemoveServiceResult,
	"changeServiceAlias": QueryChangeServiceAlias, "changeServiceAliasResult": QueryChangeServiceAliasResult,
	"serviceStatus": QueryServiceStatus, "serviceStatusResult": QueryServiceStatusResult,
	"listServices": QueryListServices, "listServicesResult": QueryListServicesResult,
	"sendTo": QuerySendToService, "sendToResult": QuerySendToServiceResult,
	"getMessage": QueryGetMessage, "getMessageResult": QueryGetMessageResult,
	"messagesList": QueryListMessages, "messagesListResult": QueryListMessagesResult,
	"eventSubscribe": QueryEventSubscribe, "eventSubscribeResult": QueryEventSubscribeResult,
}
var queryNameMap map[QueryType]string

func init() {
	queryNameMap = make(map[QueryType]string)
	for k, v := range queryTypeMap {
		queryNameMap[v] = k
	}
}

// Query struct
type Query struct {
	Type    QueryType
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
	case QueryGetConfigResult:
		var p Config
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case QueryProtocolListResult:
		var p ProtocolListResult
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
		/*
			case QueryProtocolInfo:
				if len(data) > 0 {
					var p handlers.ProtocolInfoFilter
					if err := json.Unmarshal(data, &p); err != nil {
						return err
					}
					c.Payload = &p
				}
			case QueryProtocolDiscover:
				var p handlers.ProtocolDiscoverInput
				if err := json.Unmarshal(data, &p); err != nil {
					return err
				}
				c.Payload = &p
			case QueryProtocolDiscovery:
				var p handlers.ProtocolDiscoveryInput
				if err := json.Unmarshal(data, &p); err != nil {
					return err
				}
				c.Payload = &p
			case QueryAddService:
				var p handlers.ServiceEntry
				if err := json.Unmarshal(data, &p); err != nil {
					return err
				}
				c.Payload = &p
			case QueryRemoveService:
				var p handlers.ServiceID
				if err := json.Unmarshal(data, &p); err != nil {
					return err
				}
				c.Payload = &p
			case QueryChangeServiceAlias:
				var p handlers.ChangeServiceAliasQuery
				if err := json.Unmarshal(data, &p); err != nil {
					return err
				}
				c.Payload = &p
			case QueryServiceStatus:
				var p handlers.ServiceID
				if err := json.Unmarshal(data, &p); err != nil {
					return err
				}
				c.Payload = &p
			case QueryListServices:
				var p handlers.ListServicesInput
				if err := json.Unmarshal(data, &p); err != nil {
					return err
				}
				c.Payload = &p
			case QuerySendToService:
				var p handlers.SendToServiceInput
				if err := json.Unmarshal(data, &p); err != nil {
					return err
				}
				c.Payload = &p
			case QueryGetMessage:
				var p uuid.UUID
				if err := json.Unmarshal(data, &p); err != nil {
					return err
				}
				c.Payload = &p
			case QueryListMessages:
				var p handlers.ListMessagesInput
				if err := json.Unmarshal(data, &p); err != nil {
					return err
				}
				c.Payload = &p
		*/
	case QueryEventSubscribe:
		var p Subscription
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	}
	return nil
}

// methods to create request queries

func NewQueryRestart(id string) *Query {
	return &Query{Type: QueryRestart, ID: id}
}

func NewQueryGetConfig(id string) *Query {
	return &Query{Type: QueryGetConfig, ID: id}
}

func NewQueryProtocolList(id string) *Query {
	return &Query{Type: QueryProtocolList, ID: id}
}

func NewQueryTransportList(id string) *Query {
	return &Query{Type: QueryTransportList, ID: id}
}

func NewQueryProtocolInfo(id string) *Query {
	return &Query{Type: QueryProtocolInfo, ID: id}
}

func NewQueryProtocolDiscover(id string) *Query {
	return &Query{Type: QueryProtocolDiscover, ID: id}
}

func NewQueryProtocolDiscovery(id string) *Query {
	return &Query{Type: QueryProtocolDiscovery, ID: id}
}

func NewQueryAddService(id string) *Query {
	return &Query{Type: QueryAddService, ID: id}
}

func NewQueryRemoveService(id string) *Query {
	return &Query{Type: QueryRemoveService, ID: id}
}

func NewQueryChangeServiceAlias(id string) *Query {
	return &Query{Type: QueryChangeServiceAlias, ID: id}
}

func NewQueryServiceStatus(id string) *Query {
	return &Query{Type: QueryServiceStatus, ID: id}
}

func NewQueryListServices(id string) *Query {
	return &Query{Type: QueryListServices, ID: id}
}

func NewQuerySendToService(id string) *Query {
	return &Query{Type: QuerySendToService, ID: id}
}

func NewQueryGetMessage(id string) *Query {
	return &Query{Type: QueryGetMessage, ID: id}
}

func NewQueryListMessages(id string) *Query {
	return &Query{Type: QueryListMessages, ID: id}
}

func NewQueryEventSubscribe(id string) *Query {
	return &Query{Type: QueryEventSubscribe, ID: id}
}
