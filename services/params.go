package services

// ParamType type of parameter type enum
type ParamType uint8

// Parameter's data types
const (
	ParamTypeInt32 ParamType = iota
	ParamTypeBool
	ParamTypeString
	ParamTypeEnum
)

// ParamInfo parameter description
type ParamInfo struct {
	Description  string
	Type         ParamType
	DefaultValue string
	EnumValues   []string
}

// Params type defines named parameter collection
type Params map[string]ParamInfo

// ParamValues type defines named parameter values
type ParamValues map[string]interface{}
