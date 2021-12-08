package config

import "os"

type MessageLogFlag byte

const (
	MLFlagIgnoreReadError = MessageLogFlag(1 << iota)
)

// var flagTypeMap = map[string]flagType{
// 	"ignore-read-error":  flagIgnoreReadError,
// 	"ignore-write-error": flagIgnoreWriteError,
// }

// func parseFlags(flags string) (flagType, error) {
// 	var result flagType
// 	var err error
// 	utils.ParseOptions(flags, func(flag string) bool {
// 		if fl, ok := flagTypeMap[flag]; ok {
// 			result |= fl
// 			return true
// 		}
// 		err = fmt.Errorf("unknown flag '%v'", flag)
// 		return false
// 	})
// 	return result, err
// }

// MessageLogConfig defines message log configuration entries
type MessageLogConfig struct {
	// maximal messages log size, in bytes. must be greater or equal to 8192. Default value is 10MB
	MaxSize uint `yaml:"maxSize,omitempty" json:"maxSize,omitempty"`
	// file where messages log will be stored. If not specified or empty the message log will not persist
	File     string         `yaml:"file,omitempty" json:"file,omitempty"`
	DirMode  os.FileMode    `yaml:"dirMode,omitempty" json:"dirMode,omitempty"`
	FileMode os.FileMode    `yaml:"fileMode,omitempty" json:"fileMode,omitempty"`
	Flags    MessageLogFlag `yaml:"flags,omitempty" json:"flags,omitempty"`
}

func (c *MessageLogConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// TODO
	return nil
}

func (c *MessageLogConfig) MarshalYAML() (interface{}, error) {
	// TODO
	return nil, nil
}
