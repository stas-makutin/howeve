package httpsrv

import (
	"context"
	"encoding/json"
	"fmt"

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
	queryProtocolDiscovery
	queryProtocolDiscoveryResult
)

var queryTypeMap = map[string]queryType{
	"restart": queryRestart, "restartResult": queryRestartResult,
	"getCfg": queryGetConfig, "getConfig": queryGetConfig, "getConfigResult": queryGetConfigResult,
	"protocols": queryProtocolList, "protocolsResult": queryProtocolListResult,
	"transports": queryTransportList, "transportsResult": queryTransportListResult,
	"protocolInfo": queryProtocolInfo, "protocolInfoResult": queryProtocolInfoResult,
	"discovery": queryProtocolDiscovery, "discoveryResult": queryProtocolDiscoveryResult,
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
		return fmt.Errorf("Unexpected query %v", env.Type)
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
		return nil, fmt.Errorf("Unknown query %v", c.Type)
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
	case queryProtocolDiscovery:
		var p handlers.ProtocolDiscoveryQuery
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
	case queryProtocolDiscovery:
		return &handlers.ProtocolDiscovery{RequestHeader: *handlers.NewRequestHeader(c.ID), ProtocolDiscoveryQuery: c.Payload.(*handlers.ProtocolDiscoveryQuery)}
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
	case *handlers.ProtocolDiscoveryResult:
		return &Query{Type: queryProtocolDiscoveryResult, ID: e.TraceID(), Payload: e.ProtocolDiscoveryQueryResult}
	}
	return nil
}
