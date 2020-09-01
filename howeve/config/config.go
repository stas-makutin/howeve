package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/stas-makutin/howeve/howeve/utils"
	"gopkg.in/yaml.v3"
)

// Error func
type Error func(msg string)

// Reader func
type Reader func(cfg *Config, cfgError Error)

// Writer func
type Writer func(cfg *Config)

var readers []Reader
var writers []Writer

var writeConfiguration func(restart bool) bool

// AddReader func
func AddReader(r Reader) {
	readers = append(readers, r)
}

// AddWriter func
func AddWriter(w Writer) {
	writers = append(writers, w)
}

func readConfig(cfgFile string) (*Config, error) {
	var config Config

	err := func() error {
		file, err := os.Open(cfgFile)
		if err != nil {
			return fmt.Errorf("unable to open configuration file: %v", err)
		}
		defer file.Close()
		err = yaml.NewDecoder(file).Decode(&config)
		if err != nil {
			return fmt.Errorf("unable to parse configuration file: %v", err)
		}
		return nil
	}()
	if err != nil {
		return nil, err
	}

	var errStr strings.Builder
	ce := func(msg string) {
		errStr.WriteString(utils.NewLine + msg)
	}

	for _, r := range readers {
		r(&config, ce)
	}
	if errStr.Len() > 0 {
		return nil, fmt.Errorf("the configuration file '%v' is invalid:%v", cfgFile, errStr.String())
	}

	return &config, nil
}

func writeConfig(restart bool) bool {
	if writeConfiguration != nil {
		return writeConfiguration(restart)
	}
	return false
}
