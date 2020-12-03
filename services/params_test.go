package services

import (
	"fmt"
	"testing"
)

func TestParamsParse(t *testing.T) {
	params := Params{
		"paramInt8":   &ParamInfo{Type: ParamTypeInt8},
		"paramInt16":  &ParamInfo{Type: ParamTypeInt16},
		"paramInt32":  &ParamInfo{Type: ParamTypeInt32},
		"paramInt64":  &ParamInfo{Type: ParamTypeInt64},
		"paramUint8":  &ParamInfo{Type: ParamTypeUint8},
		"paramUint16": &ParamInfo{Type: ParamTypeUint16},
		"paramUint32": &ParamInfo{Type: ParamTypeUint32},
		"paramUint64": &ParamInfo{Type: ParamTypeUint64},
		"paramBool":   &ParamInfo{Type: ParamTypeBool},
		"paramString": &ParamInfo{Type: ParamTypeString},
		"paramEnum":   &ParamInfo{Type: ParamTypeEnum, EnumValues: []string{"enum1", "enum2"}},
	}

	tests := []struct {
		name   string
		input  string
		output interface{}
		err    error
	}{
		{"paramUnknown", "1234", nil, ErrUnknownParamName},
		{"paramInt8", "-0x6f", int8(-111), nil},
		{"paramInt8", "-0xff", nil, ErrInvalidParamValue},
		{"paramInt16", "-0x6ff", int16(-1791), nil},
		{"paramInt32", "-0x6ffffff", int32(-117440511), nil},
		{"paramInt64", "-0x6fffffffffffffff", int64(-8070450532247928831), nil},
		{"paramUint8", "0b11111111", uint8(0xff), nil},
		{"paramUint8", "-045", nil, ErrInvalidParamValue},
		{"paramUint16", "0xffff", uint16(0xffff), nil},
		{"paramUint32", "0xffffffff", uint32(0xffffffff), nil},
		{"paramUint64", "0xffffffffffffffff", uint64(0xffffffffffffffff), nil},
		{"paramBool", "0", false, nil},
		{"paramBool", "FalSe", false, nil},
		{"paramBool", "1", true, nil},
		{"paramBool", "TruE", true, nil},
		{"paramBool", "123", nil, ErrInvalidParamValue},
		{"paramString", "", "", nil},
		{"paramString", "1234", "1234", nil},
		{"paramEnum", "enum1", "enum1", nil},
		{"paramEnum", "enum2", "enum2", nil},
		{"paramEnum", "enum3", nil, ErrInvalidParamValue},
		{"paramEnum", "", nil, ErrInvalidParamValue},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Test %d (%s)", i, test.name), func(t *testing.T) {
			output, err := params.Parse(test.name, test.input)
			if output != test.output || err != test.err {
				t.Errorf("expected %v, %v; got %v %v", test.output, test.err, output, err)
			}
		})
	}
}
