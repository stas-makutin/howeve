package serial

import "github.com/stas-makutin/howeve/defs"

// serial transport parameters names
const (
	ParamNameBaudRate     = "baudRate"
	ParamNameDataBits     = "dataBits"
	ParamNameParity       = "parity"
	ParamNameStopBits     = "stopBits"
	ParamNameWriteTimeout = "writeTimeout"
)

var TransportInfo *defs.TransportInfo = &defs.TransportInfo{
	Name: "Serial",
	Params: defs.Params{
		ParamNameBaudRate: {
			Description:  "The serial port bitrate",
			Type:         defs.ParamTypeInt32,
			DefaultValue: "115200",
		},
		ParamNameDataBits: {
			Description:  "The size of the character, bits",
			Type:         defs.ParamTypeEnum,
			DefaultValue: "8",
			EnumValues:   []string{"5", "6", "7", "8"},
		},
		ParamNameParity: {
			Description:  "The parity",
			Type:         defs.ParamTypeEnum,
			DefaultValue: "none",
			EnumValues:   []string{"none", "odd", "even", "mark", "space"},
		},
		ParamNameStopBits: {
			Description:  "The number of stop bits",
			Type:         defs.ParamTypeEnum,
			DefaultValue: "1",
			EnumValues:   []string{"1", "1.5", "2"},
		},
		ParamNameWriteTimeout: {
			Description:  "The write timeout, millisecons",
			Type:         defs.ParamTypeUint32,
			DefaultValue: "3000",
		},
	},
}
