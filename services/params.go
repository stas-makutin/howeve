package services

import (
	"errors"
	"strconv"
	"strings"
)

// ParamType type of parameter type enum
type ParamType uint8

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

func (pt ParamType) String() string {
	switch pt {
	case ParamTypeInt8:
		return "int8"
	case ParamTypeInt16:
		return "int16"
	case ParamTypeInt32:
		return "int32"
	case ParamTypeInt64:
		return "int64"
	case ParamTypeUint8:
		return "uint8"
	case ParamTypeUint16:
		return "uint16"
	case ParamTypeUint32:
		return "uint32"
	case ParamTypeUint64:
		return "uint64"
	case ParamTypeBool:
		return "bool"
	case ParamTypeEnum:
		return "enum"
	}
	return "string"
}

// ParamInfo parameter description
type ParamInfo struct {
	Description  string
	Type         ParamType
	DefaultValue string
	EnumValues   []string
}

// Params type defines named parameter collection
type Params map[string]*ParamInfo

// ParamValues type defines named parameter values
type ParamValues map[string]interface{}

// ErrUnknownParamName is the error for unknown parameter name
var ErrUnknownParamName error = errors.New("The parameter name is unknown")

// ErrInvalidParamValue is the error for not valid parameter value
var ErrInvalidParamValue error = errors.New("The parameter value is not valid")

// Parse function validates parameter name and parses its value
func (p Params) Parse(name, value string) (interface{}, error) {
	info, ok := p[name]
	if !ok {
		return nil, ErrUnknownParamName
	}
	switch info.Type {
	case ParamTypeInt8, ParamTypeInt16, ParamTypeInt32, ParamTypeInt64:
		bitSize := 64
		if ParamTypeInt8 == info.Type {
			bitSize = 8
		} else if ParamTypeInt16 == info.Type {
			bitSize = 16
		} else if ParamTypeInt32 == info.Type {
			bitSize = 32
		}
		if v, err := strconv.ParseInt(value, 0, bitSize); err == nil {
			if ParamTypeInt8 == info.Type {
				return int8(v), nil
			} else if ParamTypeInt16 == info.Type {
				return int16(v), nil
			} else if ParamTypeInt32 == info.Type {
				return int32(v), nil
			}
			return v, nil
		}
	case ParamTypeUint8, ParamTypeUint16, ParamTypeUint32, ParamTypeUint64:
		bitSize := 64
		if ParamTypeUint8 == info.Type {
			bitSize = 8
		} else if ParamTypeUint16 == info.Type {
			bitSize = 16
		} else if ParamTypeUint32 == info.Type {
			bitSize = 32
		}
		if v, err := strconv.ParseUint(value, 0, bitSize); err == nil {
			if ParamTypeUint8 == info.Type {
				return uint8(v), nil
			} else if ParamTypeUint16 == info.Type {
				return uint16(v), nil
			} else if ParamTypeUint32 == info.Type {
				return uint32(v), nil
			}
			return v, nil
		}
	case ParamTypeBool:
		value = strings.ToLower(value)
		if "1" == value || "true" == value {
			return true, nil
		}
		if "0" == value || "false" == value {
			return false, nil
		}
	case ParamTypeString:
		return value, nil
	case ParamTypeEnum:
		for _, v := range info.EnumValues {
			if v == value {
				return value, nil
			}
		}
	}
	return nil, ErrInvalidParamValue
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
