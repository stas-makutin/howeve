package services

import (
	"strings"

	"github.com/stas-makutin/howeve/defs"
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
var Transports = map[defs.TransportIdentifier]*defs.TransportInfo{
	defs.TransportSerial: {
		Name: "Serial",
		Params: defs.Params{
			ParamNameSerialBaudRate: {
				Description:  "The serial port bitrate",
				Type:         defs.ParamTypeInt32,
				DefaultValue: "115200",
			},
			ParamNameSerialDataBits: {
				Description:  "The size of the character, bits",
				Type:         defs.ParamTypeEnum,
				DefaultValue: "8",
				EnumValues:   []string{"5", "6", "7", "8"},
			},
			ParamNameSerialParity: {
				Description:  "The parity",
				Type:         defs.ParamTypeEnum,
				DefaultValue: "none",
				EnumValues:   []string{"none", "odd", "even", "mark", "space"},
			},
			ParamNameSerialStopBits: {
				Description:  "The number of stop bits",
				Type:         defs.ParamTypeEnum,
				DefaultValue: "1",
				EnumValues:   []string{"1", "1.5", "2"},
			},
			ParamNameSerialWriteTimeout: {
				Description:  "The write timeout, millisecons",
				Type:         defs.ParamTypeUint32,
				DefaultValue: "3000",
			},
		},
	},
}

// Protocols contains protocols definitions
var Protocols = map[defs.ProtocolIdentifier]*defs.ProtocolInfo{
	defs.ProtocolZWave: {
		Name: "Z-Wave",
		Transports: map[defs.TransportIdentifier]*defs.ProtocolTransportOptions{
			defs.TransportSerial: {
				DiscoveryFunc: zwave.DiscoverySerial,
				Params: defs.Params{
					ParamNameSerialDataBits: &defs.ParamInfo{
						Type:         defs.ParamTypeString,
						DefaultValue: "8",
						Flags:        defs.ParamFlagConst,
					},
					ParamNameSerialWriteTimeout: &defs.ParamInfo{
						Type:         defs.ParamTypeUint32,
						DefaultValue: "3000",
						Flags:        defs.ParamFlagConst,
					},
				},
			},
		},
	},
}

// TransportByName resolves transport name into identifier
func TransportByName(name string) (defs.TransportIdentifier, bool) {
	for id, ti := range Transports {
		if strings.EqualFold(name, ti.Name) {
			return id, true
		}
	}
	return 0, false
}

// ProtocolByName resolves protocol name into identifier
func ProtocolByName(name string) (defs.ProtocolIdentifier, bool) {
	for id, pi := range Protocols {
		if strings.EqualFold(name, pi.Name) {
			return id, true
		}
	}
	return 0, false
}
