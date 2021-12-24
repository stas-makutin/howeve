package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stas-makutin/howeve/utils"
	"gopkg.in/yaml.v3"
)

type MessageLogFlag byte

const (
	MLFlagIgnoreReadError = MessageLogFlag(1 << iota)
)

var messageLogFlagMap = map[string]MessageLogFlag{
	"ignore-read-error": MLFlagIgnoreReadError,
}

func (flags *MessageLogFlag) UnmarshalYAML(node *yaml.Node) (err error) {
	*flags = 0
	utils.ParseOptions(node.Value, func(flag string) bool {
		if fl, ok := messageLogFlagMap[flag]; ok {
			*flags |= fl
			return true
		}
		err = fmt.Errorf("line %d, column %d: unknown flag '%v'", node.Line, node.Column, flag)
		return false
	})
	return err
}

func (flags MessageLogFlag) String() string {
	var res strings.Builder
	for s, mask := range messageLogFlagMap {
		if (flags & mask) != 0 {
			if res.Len() > 0 {
				res.WriteString(",")
			}
			res.WriteString(s)
		}
	}
	return res.String()
}

func (flags MessageLogFlag) MarshalYAML() (interface{}, error) {
	return flags.String(), nil
}

func (flags MessageLogFlag) MarshalJSON() ([]byte, error) {
	return json.Marshal(flags.String())
}

// MessageLogConfig defines message log configuration entries
type MessageLogConfig struct {
	// maximal messages log size, in bytes. must be greater or equal to 8192. Default value is 10MB
	MaxSize SizeType `yaml:"maxSize,omitempty" json:"maxSize,omitempty"`
	// file where messages log will be stored. If not specified or empty the message log will not persist
	File       string         `yaml:"file,omitempty" json:"file,omitempty"`
	DirMode    FileMode       `yaml:"dirMode,omitempty" json:"dirMode,omitempty"`
	FileMode   FileMode       `yaml:"fileMode,omitempty" json:"fileMode,omitempty"`
	Flags      MessageLogFlag `yaml:"flags,omitempty" json:"flags,omitempty"`
	AutoPesist DurationType   `yaml:"autoPersist,omitempty" json:"autoPersist,omitempty"`
}
