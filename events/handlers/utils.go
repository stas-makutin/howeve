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

func validateServiceKey(key *defs.ServiceKey) *ErrorInfo {
	if key != nil {
		_, _, err := defs.ResolveProtocolAndTransport(key.Protocol, key.Transport)
		return handleProtocolErrors(err, key.Protocol, key.Transport)
	}
	return newErrorInfo(ErrorServiceNoKey, nil)
}

func validateServiceID(key *defs.ServiceKey, alias string) *ErrorInfo {
	if key != nil {
		return validateServiceKey(key)
	}
	if alias == "" {
		return newErrorInfo(ErrorServiceNoID, nil)
	}
	return nil
}
