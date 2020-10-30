package httpsrv

import (
	"encoding/json"
	"fmt"

	"github.com/stas-makutin/howeve/events/handlers"
)

type queryType uint16

const (
	queryRestart = queryType(iota)
	queryRestartResult
	queryGetConfig
	queryGetConfigResult
)

var queryTypeMap = map[string]queryType{
	"restart": queryRestart,
	"getCfg":  queryGetConfig, "getConfig": queryGetConfig,
}

var queryNameMap = map[queryType]string{
	queryRestart:         "restart",
	queryRestartResult:   "restartResult",
	queryGetConfig:       "getConfig",
	queryGetConfigResult: "getConfigResult",
}

// Query struct
type Query struct {
	Type    queryType   `json:"q"`
	ID      string      `json:"i"`
	Payload interface{} `json:"p,omitempty"`
}

// UnmarshalJSON func
func (c *Query) UnmarshalJSON(data []byte) error {
	var env struct {
		Type    string          `json:"q"`
		ID      string          `json:"i"`
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
		ID      string      `json:"i"`
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
		return handlers.Restart{RequestHeader: *handlers.NewRequestHeader(c.ID)}
	case queryGetConfig:
		return handlers.ConfigGet{RequestHeader: *handlers.NewRequestHeader(c.ID)}
	}
	return nil
}

func queryFromEvent(event interface{}) *Query {
	switch e := event.(type) {
	case *handlers.RestartResult:
		return &Query{Type: queryRestartResult, ID: e.ID}
	case *handlers.ConfigGetResult:
		return &Query{Type: queryGetConfigResult, ID: e.ID, Payload: e.Config}
	}
	return nil
}
