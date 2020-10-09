package httpsrv

import (
	"encoding/json"
	"fmt"
)

type queryType uint16

const (
	queryRestart = queryType(iota)
	queryGetConfig
)

// WebSocketQuery struct
type WebSocketQuery struct {
	Type    queryType   `json:"q"`
	ID      string      `json:"i"`
	Payload interface{} `json:"p,omitempty"`
}

// UnmarshalJSON func
func (c *WebSocketQuery) UnmarshalJSON(data []byte) error {
	var env struct {
		Type    string          `json:"q"`
		ID      string          `json:"i"`
		Payload json.RawMessage `json:"p,omitempty"`
	}
	err := json.Unmarshal(data, &env)
	if err != nil {
		return err
	}
	t, ok := map[string]queryType{
		"restart": queryRestart, "getCfg": queryGetConfig, "getConfig": queryGetConfig,
	}[env.Type]
	if !ok {
		return fmt.Errorf("Unknown query %v", env.Type)
	}
	c.Type = t
	c.ID = env.ID
	c.Payload = nil
	/*
		switch c.QueryID {
		}
	*/
	return err
}
