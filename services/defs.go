package services

import (
	"github.com/stas-makutin/howeve/services/servicedef"
	"github.com/stas-makutin/howeve/services/zwave"
)

// Transports contains transports definitions
var Transports = map[servicedef.TransportIdentifier]*servicedef.TransportInfo{
	servicedef.TransportSerial: {
		Name: "Serial",
		Params: servicedef.Params{
			"baudRate": {
				Description:  "The serial port bitrate",
				Type:         servicedef.ParamTypeInt32,
				DefaultValue: "115200",
			},
			"dataBits": {
				Description:  "The size of the character, bits",
				Type:         servicedef.ParamTypeEnum,
				DefaultValue: "8",
				EnumValues:   []string{"5", "6", "7", "8"},
			},
			"parity": {
				Description:  "The parity",
				Type:         servicedef.ParamTypeEnum,
				DefaultValue: "none",
				EnumValues:   []string{"none", "odd", "even", "mark", "space"},
			},
			"stopBits": {
				Description:  "The number of stop bits",
				Type:         servicedef.ParamTypeEnum,
				DefaultValue: "1",
				EnumValues:   []string{"1", "1.5", "2"},
			},
		},
	},
}

// Protocols contains protocols definitions
var Protocols = map[servicedef.ProtocolIdentifier]*servicedef.ProtocolInfo{
	servicedef.ProtocolZWave: {
		Name: "Z-Wave",
		Transports: map[servicedef.TransportIdentifier]*servicedef.ProtocolTransportOptions{
			servicedef.TransportSerial: {
				DiscoveryFunc: zwave.DiscoverySerial,
				Params: servicedef.Params{
					"dataBits": &servicedef.ParamInfo{
						Type:         servicedef.ParamTypeString,
						DefaultValue: "8",
						Const:        true,
					},
				},
			},
		},
	},
}
