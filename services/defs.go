package services

import (
	"github.com/stas-makutin/howeve/services/servicedef"
	"github.com/stas-makutin/howeve/services/zwave"
)

// serial transport parameters names
const (
	ParamNameSerialBaudRate     = "baudRate"
	ParamNameSerialDataBits     = "dataBits"
	ParamNameSerialParity       = "parity"
	ParamNameSerialStopBits     = "stopBits"
	ParamNameSerialWriteTimeout = "writeTimeout"
)

// Transports contains transports definitions
var Transports = map[servicedef.TransportIdentifier]*servicedef.TransportInfo{
	servicedef.TransportSerial: {
		Name: "Serial",
		Params: servicedef.Params{
			ParamNameSerialBaudRate: {
				Description:  "The serial port bitrate",
				Type:         servicedef.ParamTypeInt32,
				DefaultValue: "115200",
			},
			ParamNameSerialDataBits: {
				Description:  "The size of the character, bits",
				Type:         servicedef.ParamTypeEnum,
				DefaultValue: "8",
				EnumValues:   []string{"5", "6", "7", "8"},
			},
			ParamNameSerialParity: {
				Description:  "The parity",
				Type:         servicedef.ParamTypeEnum,
				DefaultValue: "none",
				EnumValues:   []string{"none", "odd", "even", "mark", "space"},
			},
			ParamNameSerialStopBits: {
				Description:  "The number of stop bits",
				Type:         servicedef.ParamTypeEnum,
				DefaultValue: "1",
				EnumValues:   []string{"1", "1.5", "2"},
			},
			ParamNameSerialWriteTimeout: {
				Description:  "The write timeout, millisecons",
				Type:         servicedef.ParamTypeUint32,
				DefaultValue: "3000",
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
					ParamNameSerialDataBits: &servicedef.ParamInfo{
						Type:         servicedef.ParamTypeString,
						DefaultValue: "8",
						Const:        true,
					},
					ParamNameSerialWriteTimeout: &servicedef.ParamInfo{
						Type:         servicedef.ParamTypeUint32,
						DefaultValue: "3000",
						Const:        true,
					},
				},
			},
		},
	},
}
