package services

import (
	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/defs"
	"github.com/stas-makutin/howeve/services/serial"
	"github.com/stas-makutin/howeve/services/zwave"
)

var transports = map[api.TransportIdentifier]*defs.TransportInfo{
	api.TransportSerial: serial.TransportInfo,
}

var protocols = map[api.ProtocolIdentifier]*defs.ProtocolInfo{
	api.ProtocolZWave: {
		Name: "Z-Wave",
		Transports: map[api.TransportIdentifier]*defs.ProtocolTransportOptions{
			api.TransportSerial: {
				ServiceFunc: func(entry string, params api.ParamValues) (defs.Service, error) {
					return zwave.NewService(&serial.Transport{}, entry, params)
				},
				DiscoveryFunc: zwave.DiscoverSerial,
				Params: defs.Params{
					serial.ParamNameDataBits: &defs.ParamInfo{
						Type:         defs.ParamTypeString,
						DefaultValue: "8",
						Flags:        defs.ParamFlagConst,
					},
					serial.ParamNameReadTimeout: &defs.ParamInfo{
						Type:         defs.ParamTypeUint32,
						DefaultValue: "0",
						Flags:        defs.ParamFlagConst,
					},
					serial.ParamNameWriteTimeout: &defs.ParamInfo{
						Type:         defs.ParamTypeUint32,
						DefaultValue: "0",
						Flags:        defs.ParamFlagConst,
					},
					defs.ParamNameOpenAttemptsInterval: {
						Description:  "The time interval between attempts to open serial port, milliseconds",
						Type:         defs.ParamTypeUint32,
						DefaultValue: "3000",
					},
					defs.ParamNameOutgoingMaxTTL: {
						Description:  "The time to live of outgoing message before it will be sent, milliseconds",
						Type:         defs.ParamTypeUint32,
						DefaultValue: "10000",
					},
				},
			},
		},
	},
}

func init() {
	defs.Transports = transports
	defs.Protocols = protocols
}
