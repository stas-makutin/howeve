package services

import (
	"strings"

	"github.com/stas-makutin/howeve/defs"
	"github.com/stas-makutin/howeve/services/serial"
	"github.com/stas-makutin/howeve/services/zwave"
)

// Transports contains transports definitions
var Transports = map[defs.TransportIdentifier]*defs.TransportInfo{
	defs.TransportSerial: serial.TransportInfo,
}

// Protocols contains protocols definitions
var Protocols = map[defs.ProtocolIdentifier]*defs.ProtocolInfo{
	defs.ProtocolZWave: {
		Name: "Z-Wave",
		Transports: map[defs.TransportIdentifier]*defs.ProtocolTransportOptions{
			defs.TransportSerial: {
				ServiceFunc: func(entry string, params defs.ParamValues) (defs.Service, error) {
					return zwave.NewService(&serial.Transport{}, entry, params)
				},
				DiscoveryFunc: zwave.DiscoverySerial,
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
						Type:         defs.ParamTypeInt32,
						DefaultValue: "3000",
					},
				},
			},
		},
	},
}

// TransportName return name of the transport for provided identifier
func TransportName(t defs.TransportIdentifier) string {
	if ti, ok := Transports[t]; ok {
		return ti.Name
	}
	return ""
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

// ProtocolName return name of the transport for provided identifier
func ProtocolName(p defs.ProtocolIdentifier) string {
	if pi, ok := Protocols[p]; ok {
		return pi.Name
	}
	return ""
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
