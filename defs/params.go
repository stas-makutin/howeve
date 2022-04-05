package defs

import (
	"strconv"
	"strings"

	"github.com/stas-makutin/howeve/api"
)

// ParamType type of parameter type enum
type ParamType uint8

// ParamFlags type of parameter bit flags
type ParamFlags uint8

// Parameter's data types
const (
	ParamTypeInt8 ParamType = iota
	ParamTypeInt16
	ParamTypeInt32
	ParamTypeInt64
	ParamTypeUint8
	ParamTypeUint16
	ParamTypeUint32
	ParamTypeUint64
	ParamTypeBool
	ParamTypeString
	ParamTypeEnum
)

// Parameter's bit flags
const (
	ParamFlagConst ParamFlags = 1 << iota
	ParamFlagRequired
)

func (pt ParamType) String() string {
	switch pt {
	case ParamTypeInt8:
		return api.ParamTypeInt8
	case ParamTypeInt16:
		return api.ParamTypeInt16
	case ParamTypeInt32:
		return api.ParamTypeInt32
	case ParamTypeInt64:
		return api.ParamTypeInt64
	case ParamTypeUint8:
		return api.ParamTypeUint8
	case ParamTypeUint16:
		return api.ParamTypeUint16
	case ParamTypeUint32:
		return api.ParamTypeUint32
	case ParamTypeUint64:
		return api.ParamTypeUint64
	case ParamTypeBool:
		return api.ParamTypeBool
	case ParamTypeEnum:
		return api.ParamTypeEnum
	}
	return api.ParamTypeString
}

// ParamInfo parameter description
type ParamInfo struct {
	Description  string
	Type         ParamType
	Flags        ParamFlags
	DefaultValue string
	EnumValues   []string
}

// Params type defines named parameter collection
type Params map[string]*ParamInfo

const (
	UnknownParamName = iota
	InvalidParamValue
	NoRequiredParam
)

// ParseError is the structure which contains parameter parse error
type ParseError struct {
	Code  int
	Name  string
	Value string
}

// NewParseError creates new parse error
func NewParseError(code int, name, value string) error {
	return &ParseError{Code: code, Name: name, Value: value}
}

// Error is the implementation of error interface
func (pe *ParseError) Error() string {
	switch pe.Code {
	case UnknownParamName:
		return "the parameter name is unknown"
	case NoRequiredParam:
		return "the required parameter is missing"
	}
	return "the parameter value is not valid"
}

func (pe *ParseError) Is(target error) bool {
	t, ok := target.(*ParseError)
	if !ok {
		return false
	}
	return *pe == *t
}

// Parse function parses provided parameter value
func (p ParamInfo) Parse(name, value string) (interface{}, error) {
	switch p.Type {
	case ParamTypeInt8, ParamTypeInt16, ParamTypeInt32, ParamTypeInt64:
		bitSize := 64
		if ParamTypeInt8 == p.Type {
			bitSize = 8
		} else if ParamTypeInt16 == p.Type {
			bitSize = 16
		} else if ParamTypeInt32 == p.Type {
			bitSize = 32
		}
		if v, err := strconv.ParseInt(value, 0, bitSize); err == nil {
			if ParamTypeInt8 == p.Type {
				return int8(v), nil
			} else if ParamTypeInt16 == p.Type {
				return int16(v), nil
			} else if ParamTypeInt32 == p.Type {
				return int32(v), nil
			}
			return v, nil
		}
	case ParamTypeUint8, ParamTypeUint16, ParamTypeUint32, ParamTypeUint64:
		bitSize := 64
		if ParamTypeUint8 == p.Type {
			bitSize = 8
		} else if ParamTypeUint16 == p.Type {
			bitSize = 16
		} else if ParamTypeUint32 == p.Type {
			bitSize = 32
		}
		if v, err := strconv.ParseUint(value, 0, bitSize); err == nil {
			if ParamTypeUint8 == p.Type {
				return uint8(v), nil
			} else if ParamTypeUint16 == p.Type {
				return uint16(v), nil
			} else if ParamTypeUint32 == p.Type {
				return uint32(v), nil
			}
			return v, nil
		}
	case ParamTypeBool:
		value = strings.ToLower(value)
		if value == "1" || value == "true" {
			return true, nil
		}
		if value == "0" || value == "false" {
			return false, nil
		}
	case ParamTypeString:
		return value, nil
	case ParamTypeEnum:
		for _, v := range p.EnumValues {
			if v == value {
				return value, nil
			}
		}
	}
	return nil, NewParseError(InvalidParamValue, name, value)
}

// Parse function validates parameter name and parses its value
func (p Params) Parse(name, value string) (interface{}, error) {
	param, ok := p[name]
	if !ok {
		return nil, NewParseError(UnknownParamName, name, value)
	}
	return param.Parse(name, value)
}

// ParseValues function validates parameters and parses their value. Returns values map or parameter name + associated error
func (p Params) ParseValues(values api.RawParamValues) (api.ParamValues, error) {
	rv := make(api.ParamValues)

	for name, param := range p {
		value, ok := values[name]
		if ok {
			if param.Flags&ParamFlagConst != 0 {
				value = param.DefaultValue
			}
		} else {
			if param.Flags&ParamFlagRequired != 0 {
				return nil, NewParseError(NoRequiredParam, name, "")
			}
			value = param.DefaultValue
		}
		if v, err := param.Parse(name, value); err == nil {
			rv[name] = v
		} else {
			return nil, err
		}
	}

	return rv, nil
}

// Merge returns copy of combined parameters with subordinate parameters
func (p Params) Merge(subp Params) (result Params) {
	result = make(Params)
	for k, v := range p {
		result[k] = v
	}
	for k, v := range subp {
		if _, ok := result[k]; !ok {
			result[k] = v
		}
	}
	return
}
