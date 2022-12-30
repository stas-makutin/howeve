package core

import (
	"fmt"
	"strconv"

	"github.com/stas-makutin/howeve/api"
)

type Parameter struct {
	Name  string
	Value string
}

type RenderParameter struct {
	Parameter
	Index int
	Data  interface{}
}

type ParameterEnumData struct {
	Options []ParameterEnumOption
}

type ParameterBoolData struct {
	BoolValue bool
}

type ParameterIntData struct {
	Minimum string
	Maximum string
	Error   string
}

type ParameterEnumOption struct {
	Name     string
	Selected bool
}

type Parameters []Parameter

func ParamDefaultValue(pi *api.ParamInfoEntry) string {
	if pi.DefaultValue != "" {
		return pi.DefaultValue
	} else if pi.Type == api.ParamTypeEnum {
		return pi.EnumValues[0]
	} else if pi.Type == api.ParamTypeBool {
		return "false"
	} else if pi.Type != api.ParamTypeString {
		return "0"
	}
	return ""
}

func (params Parameters) FirstAvailable(transport *api.ProtocolTransportInfoEntry) (string, *api.ParamInfoEntry) {
ParamLoop:
	for name, pi := range transport.Params {
		for _, p := range params {
			if name == p.Name {
				continue ParamLoop
			}
		}
		return name, pi
	}
	return "", nil
}

func (params Parameters) AvailableNames(transport *api.ProtocolTransportInfoEntry) []string {
	var names []string
ParamLoop:
	for name := range transport.Params {
		for _, p := range params {
			if name == p.Name {
				continue ParamLoop
			}
		}
		names = append(names, name)
	}
	return names
}

func (params Parameters) AppendAvailable(transport *api.ProtocolTransportInfoEntry) (Parameters, bool) {
	name, pi := params.FirstAvailable(transport)
	if pi == nil {
		return params, false
	}
	return append(params, Parameter{Name: name, Value: ParamDefaultValue(pi)}), true
}

func (params Parameters) Replace(index int, name string, transport *api.ProtocolTransportInfoEntry) bool {
	if index < 0 || index >= len(params) {
		return false
	}
	param := &(params[index])
	if param.Name == name {
		return false
	}
	pi, ok := transport.Params[name]
	if !ok {
		return false
	}
	param.Name = name
	param.Value = ParamDefaultValue(pi)
	return true
}

func (params Parameters) ChangeValue(index int, value string) bool {
	if index < 0 || index >= len(params) {
		return false
	}
	param := &(params[index])
	if param.Value == value {
		return false
	}
	param.Value = value
	return true
}

func (params Parameters) Remove(index int) (Parameters, bool) {
	if index < 0 || index >= len(params) {
		return params, false
	}
	return append(params[:index], params[index+1:]...), true
}

func (params Parameters) ToRender(transport *api.ProtocolTransportInfoEntry) ([]*RenderParameter, bool) {
	result := make([]*RenderParameter, 0, len(params))
	valid := true
	for i, p := range params {
		pi, ok := transport.Params[p.Name]
		if !ok {
			continue
		}

		param := &RenderParameter{Parameter: p, Index: i}
		result = append(result, param)
		switch pi.Type {
		case api.ParamTypeEnum:
			data := &ParameterEnumData{}
			for _, enumValue := range pi.EnumValues {
				data.Options = append(data.Options, ParameterEnumOption{Name: enumValue, Selected: enumValue == p.Value})
			}
			param.Data = data
		case api.ParamTypeBool:
			param.Data = &ParameterBoolData{BoolValue: p.Value == "true" || p.Value == "1"}
		case api.ParamTypeString:
			// do nothing
		default: // integer parameter
			data := makeParameterIntData(param.Value, pi)
			if data.Error != "" {
				valid = false
			}
			param.Data = data
		}
	}
	return result, valid
}

func makeParameterIntData(value string, pi *api.ParamInfoEntry) *ParameterIntData {
	switch pi.Type {
	case api.ParamTypeInt8:
		return newParameterIntData(
			value,
			int64(api.ParamTypeInt8Min), int64(api.ParamTypeInt8Max),
			func(v string) (int64, error) { return strconv.ParseInt(v, 10, 8) },
			func(v int64) string { return strconv.FormatInt(v, 10) },
		)
	case api.ParamTypeInt16:
		return newParameterIntData(
			value,
			int64(api.ParamTypeInt16Min), int64(api.ParamTypeInt16Max),
			func(v string) (int64, error) { return strconv.ParseInt(v, 10, 16) },
			func(v int64) string { return strconv.FormatInt(v, 10) },
		)
	case api.ParamTypeInt32:
		return newParameterIntData(
			value,
			int64(api.ParamTypeInt32Min), int64(api.ParamTypeInt32Max),
			func(v string) (int64, error) { return strconv.ParseInt(v, 10, 32) },
			func(v int64) string { return strconv.FormatInt(v, 10) },
		)
	case api.ParamTypeInt64:
		return newParameterIntData(
			value,
			int64(api.ParamTypeInt64Min), int64(api.ParamTypeInt64Max),
			func(v string) (int64, error) { return strconv.ParseInt(v, 10, 64) },
			func(v int64) string { return strconv.FormatInt(v, 10) },
		)
	case api.ParamTypeUint16:
		return newParameterIntData(
			value,
			uint64(api.ParamTypeUint16Min), uint64(api.ParamTypeUint16Max),
			func(v string) (uint64, error) { return strconv.ParseUint(v, 10, 16) },
			func(v uint64) string { return strconv.FormatUint(v, 10) },
		)
	case api.ParamTypeUint32:
		return newParameterIntData(
			value,
			uint64(api.ParamTypeUint32Min), uint64(api.ParamTypeUint32Max),
			func(v string) (uint64, error) { return strconv.ParseUint(v, 10, 32) },
			func(v uint64) string { return strconv.FormatUint(v, 10) },
		)
	case api.ParamTypeUint64:
		return newParameterIntData(
			value,
			uint64(api.ParamTypeUint64Min), uint64(api.ParamTypeUint64Max),
			func(v string) (uint64, error) { return strconv.ParseUint(v, 10, 64) },
			func(v uint64) string { return strconv.FormatUint(v, 10) },
		)
	default:
		return newParameterIntData(
			value,
			uint64(api.ParamTypeUint8Min), uint64(api.ParamTypeUint8Max),
			func(v string) (uint64, error) { return strconv.ParseUint(v, 10, 8) },
			func(v uint64) string { return strconv.FormatUint(v, 10) },
		)
	}
}

func newParameterIntData[T int64 | uint64](value string, min, max T, parse func(string) (T, error), format func(T) string) *ParameterIntData {
	data := &ParameterIntData{}
	data.Minimum = format(min)
	data.Maximum = format(max)
	v, err := parse(value)
	if err != nil {
		data.Error = fmt.Sprintf("The value must be equal or greater than %s and equal or less than %s", data.Minimum, data.Maximum)
	} else if v < min {
		data.Error = fmt.Sprintf("The value must be equal or greater than %s", data.Minimum)
	} else if v > max {
		data.Error = fmt.Sprintf("The value must be equal or less than %s", data.Maximum)
	}
	return data
}
