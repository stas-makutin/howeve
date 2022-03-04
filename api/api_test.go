package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSerialization(t *testing.T) {
	uid := uuid.New()
	quieries := []*Query{
		NewQueryRestart("qr"),
		{Type: QueryRestartResult, ID: "qrr"},

		NewQueryGetConfig("qgc"),
		{
			Type: QueryGetConfigResult, ID: "qrgc",
			Payload: &Config{
				WorkingDirectory: "test",
			},
		},

		NewQueryProtocolList("gpl"),
		{
			Type: QueryProtocolListResult, ID: "qrpl",
			Payload: &ProtocolListResult{
				Protocols: []*ProtocolListEntry{
					{ProtocolZWave, "Protocol 1"},
					{ProtocolZWave, "Protocol 2"},
				},
			},
		},

		NewQueryTransportList("qtl"),
		{
			Type: QueryTransportListResult, ID: "qrtl",
			Payload: &TransportListResult{
				Transports: []*TransportListEntry{
					{TransportSerial, "Transport 1"},
					{TransportSerial, "Transport 2"},
				},
			},
		},

		NewQueryProtocolInfo("qpi1", nil),
		NewQueryProtocolInfo("qpi2", &ProtocolInfo{
			Protocols:  []ProtocolIdentifier{ProtocolZWave, ProtocolZWave},
			Transports: []TransportIdentifier{TransportSerial, TransportSerial},
		}),
		{
			Type: QueryProtocolInfoResult, ID: "qrpi",
			Payload: &ProtocolInfoResult{
				Protocols: []*ProtocolInfoEntry{
					{ProtocolZWave, true, "Protocol 1", []*ProtocolTransportInfoEntry{
						{TransportSerial, true, "Transport 1", map[string]*ParamInfoEntry{
							"Param 1": {"Description 1", "int8", "", nil},
							"Param 2": {"Description 2", "string", "default", nil},
						}, true, map[string]*ParamInfoEntry{
							"Param A": {"Description A", "bool", "", nil},
							"Param B": {"Description B", "enum", "one", []string{"one", "two"}},
						}},
					}},
					{ProtocolZWave, false, "Protocol 2", []*ProtocolTransportInfoEntry{
						{TransportSerial, false, "Transport A", map[string]*ParamInfoEntry{
							"Param 1": {"Description 1", "int16", "", nil},
							"Param 2": {"Description 2", "int32", "12345", nil},
						}, true, map[string]*ParamInfoEntry{
							"Param A": {"Description A", "int64", "", nil},
							"Param B": {"Description B", "uint8", "255", nil},
						}},
						{TransportSerial, false, "Transport B", map[string]*ParamInfoEntry{
							"Param 1": {"Description 1", "uint16", "345", nil},
							"Param 2": {"Description 2", "uint32", "A", nil},
						}, true, map[string]*ParamInfoEntry{
							"Param A": {"Description A", "uint64", "BC", nil},
						}},
					}},
				},
			},
		},

		NewQueryProtocolDiscover("qpd", &ProtocolDiscover{
			Protocol: ProtocolZWave, Transport: TransportSerial, Params: RawParamValues{"param1": "value1", "param2": "value2"},
		}),
		{
			Type: QueryProtocolDiscoverResult, ID: "qrpd", Payload: ProtocolDiscoverResult{
				ID:    &uid,
				Error: &ErrorInfo{ErrorDiscoveryFailed, "message", nil, nil},
			},
		},

		NewQueryProtocolDiscovery("qpdy", &ProtocolDiscovery{ID: uuid.New(), Stop: true}),
		{
			Type: QueryProtocolDiscoveryResult, ID: "qrpdy", Payload: &ProtocolDiscoveryResult{
				ID: uuid.New(), Error: &ErrorInfo{ErrorDiscoveryFailed, "message", []interface{}{1, "aaa", 1.3}, fmt.Errorf("Some error")},
				Entries: []*DiscoveryEntry{
					{
						ServiceKey: ServiceKey{ProtocolZWave, TransportSerial, "COM1"}, Description: "Some description",
						ParamValues: ParamValues{"a": 123, "b": "str", "c": true},
					},
					{
						ServiceKey: ServiceKey{ProtocolZWave, TransportSerial, "COM3"}, Description: "Another description",
						ParamValues: ParamValues{"abc": nil, "def": 12.23},
					},
				},
			},
		},

		{
			Type: QueryProtocolDiscoveryStarted, ID: "qpds", Payload: &ProtocolDiscoveryStarted{
				ProtocolDiscover: ProtocolDiscover{ProtocolZWave, TransportSerial, RawParamValues{"p1": "v1", "p2": "v2", "p3": "v3"}},
				ID:               uuid.New(),
			},
		},
		{
			Type: QueryProtocolDiscoveryFinished, ID: "qpdf", Payload: &ProtocolDiscoveryResult{
				ID: uuid.New(), Error: &ErrorInfo{ErrorDiscoveryBusy, "message", nil, nil},
				Entries: []*DiscoveryEntry{
					{
						ServiceKey: ServiceKey{ProtocolZWave, TransportSerial, "COM1"}, Description: "Some description",
						ParamValues: ParamValues{"a": 123, "b": "str", "c": true},
					},
					{
						ServiceKey: ServiceKey{ProtocolZWave, TransportSerial, "COM3"}, Description: "Another description",
						ParamValues: ParamValues{"abc": nil, "def": 12.23},
					},
				},
			},
		},

		NewQueryAddService("qas", &ServiceEntry{
			&ServiceKey{ProtocolZWave, TransportSerial, "COM3"}, RawParamValues{"p1": "v1", "p2": "v2"}, "Some alias",
		}),
		{
			Type: QueryAddServiceResult, ID: "qras", Payload: &StatusReply{nil, true},
		},

		NewQueryRemoveService("qrs", &ServiceID{&ServiceKey{ProtocolZWave, TransportSerial, ""}, ""}),
		{
			Type: QueryRemoveServiceResult, ID: "qrrs", Payload: &StatusReply{
				&ErrorInfo{ErrorServiceKeyNotExists, "message", nil, nil},
				false,
			},
		},

		NewQueryChangeServiceAlias("qcsa", &ChangeServiceAlias{
			&ServiceID{&ServiceKey{ProtocolZWave, TransportSerial, "COM5"}, ""}, "Some alias",
		}),
		{
			Type: QueryChangeServiceAliasResult, ID: "qrcsa", Payload: &StatusReply{nil, true},
		},

		NewQueryServiceStatus("qss", &ServiceID{nil, "Some alias"}),
		{
			Type: QueryServiceStatusResult, ID: "qrss", Payload: &StatusReply{nil, true},
		},

		NewQueryListServices("qls", &ListServices{
			Protocols:  []ProtocolIdentifier{ProtocolZWave},
			Transports: []TransportIdentifier{TransportSerial, TransportSerial},
			Entries:    []string{"COM1", "COM2"},
			Aliases:    []string{"Alias 1", "Alias 2"},
		}),
		{
			Type: QueryListServicesResult, ID: "qrls", Payload: &ListServicesResult{
				Services: []ListServicesEntry{
					{&ServiceID{&ServiceKey{ProtocolZWave, TransportSerial, "COM5"}, ""}, &StatusReply{nil, true}},
					{&ServiceID{&ServiceKey{ProtocolZWave, TransportSerial, "COM1"}, "Alias"}, &StatusReply{nil, false}},
				},
			},
		},

		NewQuerySendToService("qsts", &SendToService{&ServiceID{nil, "Alias"}, []byte{0xab, 0xf5, 0xc6, 0xe4, 0xb8}}),
		{
			Type: QuerySendToServiceResult, ID: "qrsts", Payload: &SendToServiceResult{
				&StatusReply{nil, true},
				&Message{time.Now(), uuid.New(), OutgoingPending, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}},
			},
		},

		NewQueryGetMessage("qgm", uuid.New()),
		{
			Type: QueryGetMessageResult, ID: "qrgm",
		},

		NewQueryListMessages("qlm", nil),
		{
			Type: QueryListMessagesResult, ID: "qrlm",
		},

		{
			Type: QueryNewMessage, ID: "qrnm",
		},
		{
			Type: QueryDropMessage, ID: "qrdm",
		},
		{
			Type: QueryUpdateMessageState, ID: "qrums",
		},

		NewQueryEventSubscribe("qes", nil),
		{
			Type: QueryEventSubscribeResult, ID: "qres",
		},
	}

	t.Run("JSON serialization", func(t *testing.T) {
		for i, query := range quieries[15:16] {
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
