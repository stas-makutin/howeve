package api

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSerialization(t *testing.T) {
	quieries := []*Query{
		NewQueryRestart("id1"),
		{Type: QueryRestartResult, ID: "id2"},
		NewQueryGetConfig("id3"),
		{
			Type: QueryGetConfigResult, ID: "id4",
			Payload: &Config{
				WorkingDirectory: "test",
			},
		},
		NewQueryProtocolList("id5"),
		{
			Type: QueryProtocolListResult, ID: "id6",
			Payload: &ProtocolListResult{
				Protocols: []*ProtocolListEntry{
					{ProtocolZWave, "Protocol 1"},
					{ProtocolZWave, "Protocol 2"},
				},
			},
		},
	}

	t.Run("JSON serialization", func(t *testing.T) {
		for i, query := range quieries {
			queryType := queryNameMap[query.Type]

			// serialize
			var writter1 strings.Builder
			if err := json.NewEncoder(&writter1).Encode(query); err != nil {
				t.Errorf("Test %d (%s) - json encode failed: %v", i, queryType, err)
				continue
			}

			// deserialize
			var decodedQuery *Query
			if err := json.NewDecoder(strings.NewReader(writter1.String())).Decode(&decodedQuery); err != nil {
				t.Errorf("Test %d (%s) - json decode failed: %v", i, queryType, err)
				continue
			}

			// serialize 2
			var writter2 strings.Builder
			if err := json.NewEncoder(&writter2).Encode(decodedQuery); err != nil {
				t.Errorf("Test %d (%s) - json encode 2 failed: %v", i, queryType, err)
				continue
			}

			if writter1.String() != writter2.String() {
				t.Errorf("Test %d (%s) - encoded jsons are different", i, queryType)
			}
		}
	})
}
