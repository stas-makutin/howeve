package httpsrv

import (
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
)

var queryTypeMap = map[string]queryType{
	"restart": queryRestart, "restartResult": queryRestartResult,
	"getCfg": queryGetConfig, "getConfig": queryGetConfig, "getConfigResult": queryGetConfigResult,
}

var queryNameMap = map[queryType]string{
	queryUnexpected:      "unexpected",
	queryRestart:         "restart",
	queryRestartResult:   "restartResult",
	queryGetConfig:       "getConfig",
	queryGetConfigResult: "getConfigResult",
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
	c.Payload = env.Payload
	return nil
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

func (c *Query) toEvent() interface{} {
	switch c.Type {
	case queryRestart:
		return &handlers.Restart{RequestHeader: *handlers.NewRequestHeader(c.ID)}
	case queryGetConfig:
		return &handlers.ConfigGet{RequestHeader: *handlers.NewRequestHeader(c.ID)}
	}
	return nil
}

func (c *Query) toTargetedRequest(receiverID events.SubscriberID) interface{} {
	event := c.toEvent()
	if te, ok := event.(events.TargetedRequest); ok {
		te.SetReceiver(receiverID)
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
	}
	return nil
}
