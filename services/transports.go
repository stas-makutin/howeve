package services

// TransportIdentifier type
type TransportIdentifier uint8

// Supported transport identifiers
const (
	TransportSerial = TransportIdentifier(iota + 1)
)

// TransportInfo transport definition structure
type TransportInfo struct {
	Name string
	Params
}

// Transports contains transports definitions
var Transports = map[TransportIdentifier]*TransportInfo{
	TransportSerial: {
		Name: "Serial",
		Params: Params{
			"baudRate": {
				Description:  "The serial port bitrate",
				Type:         ParamTypeInt32,
				DefaultValue: "115200",
			},
			"dataBits": {
				Description:  "The size of the character, bits",
				Type:         ParamTypeEnum,
				DefaultValue: "8",
				EnumValues:   []string{"5", "6", "7", "8"},
			},
			"parity": {
				Description:  "The parity",
				Type:         ParamTypeEnum,
				DefaultValue: "none",
				EnumValues:   []string{"none", "odd", "even", "mark", "space"},
			},
			"stopBits": {
				Description:  "The number of stop bits",
				Type:         ParamTypeEnum,
				DefaultValue: "1",
				EnumValues:   []string{"1", "1.5", "2"},
			},
		},
	},
}
