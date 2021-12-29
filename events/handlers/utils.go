package handlers

import (
	"github.com/stas-makutin/howeve/defs"
)

func buildParamsInfo(p defs.Params) (pie map[string]*ParamInfoEntry) {
	pie = make(map[string]*ParamInfoEntry)
	for name, pi := range p {
		if pi.Flags&defs.ParamFlagConst == 0 {
			pie[name] = &ParamInfoEntry{
				Description:  pi.Description,
				Type:         pi.Type.String(),
				DefaultValue: pi.DefaultValue,
				EnumValues:   pi.EnumValues,
			}
		}
	}
	return
}

func makeServiceKey(protocol defs.ProtocolIdentifier, transport defs.TransportIdentifier, entry string) (*defs.ServiceKey, *ErrorInfo) {
	_, _, err := defs.ResolveProtocolAndTransport(protocol, transport)
	if errorInfo := handleProtocolErrors(err, protocol, transport); errorInfo != nil {
		return nil, errorInfo
	}
	return &defs.ServiceKey{Protocol: protocol, Transport: transport, Entry: entry}, nil
}
