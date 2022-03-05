package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/google/uuid"
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
	"discoveryStarted": QueryProtocolDiscoveryStarted, "discoveryFinished": QueryProtocolDiscoveryFinished,
	"addService": QueryAddService, "addServiceResult": QueryAddServiceResult,
	"removeService": QueryRemoveService, "removeServiceResult": QueryRemoveServiceResult,
	"changeServiceAlias": QueryChangeServiceAlias, "changeServiceAliasResult": QueryChangeServiceAliasResult,
	"serviceStatus": QueryServiceStatus, "serviceStatusResult": QueryServiceStatusResult,
	"listServices": QueryListServices, "listServicesResult": QueryListServicesResult,
	"sendTo": QuerySendToService, "sendToResult": QuerySendToServiceResult,
	"getMessage": QueryGetMessage, "getMessageResult": QueryGetMessageResult,
	"messagesList": QueryListMessages, "messagesListResult": QueryListMessagesResult,
	"newMessage": QueryNewMessage, "dropMessage": QueryDropMessage, "updateMessageState": QueryUpdateMessageState,
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
	if bytes.Equal(env.Payload, []byte("null")) {
		return nil
	}
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
	if c.Payload != nil {
		v := reflect.ValueOf(c.Payload)
		if !(v.Kind() == reflect.Ptr && v.IsNil()) {
			env.Payload = c.Payload
		}
	}
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
	case QueryTransportListResult:
		var p TransportListResult
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case QueryProtocolInfo:
		if len(data) > 0 {
			var p ProtocolInfo
			if err := json.Unmarshal(data, &p); err != nil {
				return err
			}
			c.Payload = &p
		}
	case QueryProtocolInfoResult:
		var p ProtocolInfoResult
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case QueryProtocolDiscover:
		var p ProtocolDiscover
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case QueryProtocolDiscoverResult:
		var p ProtocolDiscoverResult
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case QueryProtocolDiscovery:
		var p ProtocolDiscovery
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case QueryProtocolDiscoveryResult, QueryProtocolDiscoveryFinished:
		var p ProtocolDiscoveryResult
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case QueryProtocolDiscoveryStarted:
		var p ProtocolDiscoveryStarted
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case QueryAddService:
		var p ServiceEntry
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case QueryAddServiceResult, QueryRemoveServiceResult, QueryChangeServiceAliasResult, QueryServiceStatusResult:
		var p StatusReply
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case QueryRemoveService, QueryServiceStatus:
		var p ServiceID
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case QueryChangeServiceAlias:
		var p ChangeServiceAlias
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case QueryListServices:
		var p ListServices
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case QueryListServicesResult:
		var p ListServicesResult
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case QuerySendToService:
		var p SendToService
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case QuerySendToServiceResult:
		var p SendToServiceResult
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
	case QueryGetMessageResult, QueryNewMessage, QueryDropMessage:
		var p MessageEntry
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case QueryListMessages:
		var p ListMessages
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case QueryListMessagesResult:
		var p ListMessagesResult
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case QueryUpdateMessageState:
		var p UpdateMessageState
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	case QueryEventSubscribe:
		var p Subscription
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		c.Payload = &p
	}
	return nil
}
