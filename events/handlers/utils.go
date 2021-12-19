package handlers

import (
	"errors"

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

func handleParamsErrors(err error) *ErrorInfo {
	var pe *defs.ParseError
	if errors.As(err, &pe) {
		switch pe.Code {
		case defs.UnknownParamName:
			return NewErrorInfo(ErrorUnknownParameter, pe.Name)
		case defs.NoRequiredParam:
			return NewErrorInfo(ErrorNoRequiredParameter, pe.Name)
		}
		return NewErrorInfo(ErrorInvalidParameterValue, pe.Value, pe.Name)
	}
	return nil
}

type protocolAndTransport struct {
	protocol  *defs.ProtocolInfo
	transport *defs.TransportInfo
	options   *defs.ProtocolTransportOptions
}

// findProtocolAndTransport finds protocol and transport combination, if any
func findProtocolAndTransport(protocol defs.ProtocolIdentifier, transport defs.TransportIdentifier) (protocolAndTransport, *ErrorInfo) {
	pat := protocolAndTransport{}
	if pi, ok := defs.Protocols[protocol]; ok {
		pat.protocol = pi
		if ti, ok := defs.Transports[transport]; ok {
			pat.transport = ti
			if pti, ok := pi.Transports[transport]; ok {
				pat.options = pti
				return pat, nil
			} else {
				return pat, NewErrorInfo(ErrorInvalidProtocolTransport, pat.protocol.Name, protocol, pat.transport.Name, transport)
			}
		} else {
			return pat, NewErrorInfo(ErrorUnknownTransport, transport)
		}
	}
	return pat, NewErrorInfo(ErrorUnknownProtocol, protocol)
}

func makeServiceKey(protocol defs.ProtocolIdentifier, transport defs.TransportIdentifier, entry string) (*defs.ServiceKey, *ErrorInfo) {
	_, errorInfo := findProtocolAndTransport(protocol, transport)
	if errorInfo != nil {
		return nil, errorInfo
	}
	return &defs.ServiceKey{Protocol: protocol, Transport: transport, Entry: entry}, nil
}
