package messages

import (
	"fmt"

	"github.com/stas-makutin/howeve/utils"
)

type flagType byte

const (
	flagIgnoreReadError = flagType(1 << iota)
	flagIgnoreWriteError
)

var flagTypeMap = map[string]flagType{
	"ignore-read-error":  flagIgnoreReadError,
	"ignore-write-error": flagIgnoreWriteError,
}

func parseFlags(flags string) (flagType, error) {
	var result flagType
	var err error
	utils.ParseOptions(flags, func(flag string) bool {
		if fl, ok := flagTypeMap[flag]; ok {
			result |= fl
			return true
		}
		err = fmt.Errorf("unknown flag '%v'", flag)
		return false
	})
	return result, err
}
