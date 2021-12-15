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

// parse function parses input parameters according their definition
func (pv ParamsValues) parse(p defs.Params) (defs.ParamValues, *ErrorInfo) {
	values, name, err := p.ParseAll(pv)
	if err != nil {
		if errors.Is(err, defs.ErrUnknownParamName) {
			return nil, NewErrorInfo(ErrorUnknownParameter, name)
		} else if errors.Is(err, defs.ErrNoRequiredParam) {
			return nil, NewErrorInfo(ErrorNoRequiredParameter, name)
		}
		value := pv[name]
		return nil, NewErrorInfo(ErrorInvalidParameterValue, value, name)
	}
	return values, nil
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

func makeServiceEntry(protocol defs.ProtocolIdentifier, transport defs.TransportIdentifier, entry string, pv ParamsValues) (*defs.ServiceEntry, *ErrorInfo) {
	pat, errorInfo := findProtocolAndTransport(protocol, transport)
	if errorInfo != nil {
		return nil, errorInfo
	}
	params, errorInfo := pv.parse(pat.options.Params)
	if errorInfo != nil {
		return nil, errorInfo
	}
	return &defs.ServiceEntry{Key: defs.ServiceKey{Protocol: protocol, Transport: transport, Entry: entry}, Params: params}, nil
}
