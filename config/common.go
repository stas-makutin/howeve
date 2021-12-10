package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/stas-makutin/howeve/utils"
	"gopkg.in/yaml.v3"
)

type FileMode os.FileMode

func (mode *FileMode) UnmarshalYAML(node *yaml.Node) (err error) {
	var v uint32
	err = node.Decode(&v)
	*mode = FileMode(v)
	return
}

func (mode FileMode) Value() os.FileMode {
	return os.FileMode(mode)
}

func (mode FileMode) WithDefault(dv os.FileMode) os.FileMode {
	if mode == 0 {
		return dv
	}
	return mode.Value()
}

func (mode FileMode) WithDirDefault() os.FileMode {
	return mode.WithDefault(0755)
}

func (mode FileMode) WithFileDefault() os.FileMode {
	return mode.WithDefault(0644)
}

func (mode FileMode) String() string {
	return "0" + strconv.FormatUint(uint64(mode), 8)
}

func (mode FileMode) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!int",
		Value: mode.String(),
	}, nil
}

func (mode FileMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(mode.Value())
}

type SizeType int64

func (sz *SizeType) UnmarshalYAML(node *yaml.Node) error {
	value, err := utils.ParseSizeString(node.Value)
	if err != nil {
		return fmt.Errorf("line %d, column %d: %w", node.Line, node.Column, err)
	}
	*sz = SizeType(value)
	return nil
}

func (sz SizeType) Value() int64 {
	return int64(sz)
}

func (sz SizeType) String() string {
	return utils.FormatSizeString(int64(sz))
}

func (sz SizeType) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: yaml.FlowStyle,
		Value: sz.String(),
	}, nil
}

func (sz SizeType) MarshalJSON() ([]byte, error) {
	return json.Marshal(sz.Value())
}

type DurationType time.Duration

func (d *DurationType) UnmarshalYAML(node *yaml.Node) error {
	value, err := utils.ParseTimeDuration(node.Value)
	if err != nil {
		return fmt.Errorf("line %d, column %d: %w", node.Line, node.Column, err)
	}
	*d = DurationType(value)
	return nil
}

func (d DurationType) Value() time.Duration {
	return time.Duration(d)
}

func (d DurationType) String() string {
	return utils.FormatTimeDuration(time.Duration(d))
}

func (d DurationType) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: yaml.FlowStyle,
		Value: d.String(),
	}, nil
}

func (d DurationType) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Value())
}
