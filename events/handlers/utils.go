package handlers

import (
	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/defs"
)

func buildParamsInfo(p defs.Params) (pie map[string]*api.ParamInfoEntry) {
	pie = make(map[string]*api.ParamInfoEntry)
	for name, pi := range p {
		if pi.Flags&defs.ParamFlagConst == 0 {
			pie[name] = &api.ParamInfoEntry{
				Description:  pi.Description,
				Type:         pi.Type.String(),
				DefaultValue: pi.DefaultValue,
				EnumValues:   pi.EnumValues,
			}
		}
	}
	return
}

func validateServiceKey(key *api.ServiceKey) *api.ErrorInfo {
	if key != nil {
		_, _, err := defs.ResolveProtocolAndTransport(key.Protocol, key.Transport)
		return handleProtocolErrors(err, key.Protocol, key.Transport)
	}
	return newErrorInfo(api.ErrorServiceNoKey, nil)
}

func validateServiceID(key *api.ServiceKey, alias string) *api.ErrorInfo {
	if key != nil {
		return validateServiceKey(key)
	}
	if alias == "" {
		return newErrorInfo(api.ErrorServiceNoID, nil)
	}
	return nil
}
