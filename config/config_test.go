package config

import (
	"encoding/json"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConfig(t *testing.T) {
	src := `
workingDir: test
log:
    dir: test
    file: test
    dirMode: 0755
    fileMode: 0644
    maxSize: 10 MiB
    maxAge: 1 h
    backups: 3
    backupDays: 14
    archive: zip
httpServer:
    port: 4444
    maxConnections: 111
    readTimeout: 5123 ms
    readHeaderTimeout: 1 s
    writeTimeout: 1 s
    idleTimeout: 3 s
    maxHeaderBytes: 10 KiB
messageLog:
    maxSize: 10
    file: test
    dirMode: 0777
    fileMode: 0777
    flags: ignore-read-error
services:
    - alias: homezw
      protocol: zwave
      transport: serial
      entry: COM3
      params:
        param1: value1
        param2: value2
`

	var config Config
	t.Run("Parse YAML", func(t *testing.T) {
		if err := yaml.NewDecoder(strings.NewReader(src)).Decode(&config); err != nil {
			t.Error(err)
		}
	})

	var written strings.Builder
	t.Run("Write YAML", func(t *testing.T) {
		if err := yaml.NewEncoder(&written).Encode(&config); err != nil {
			t.Error(err)
		}
	})

	t.Run("Compare YAML", func(t *testing.T) {
		if strings.TrimSpace(written.String()) != strings.TrimSpace(src) {
			t.Error("written and source YAML are different")
		}
	})

	var writtenJson strings.Builder
	t.Run("Write JSON", func(t *testing.T) {
		if err := json.NewEncoder(&writtenJson).Encode(&config); err != nil {
			t.Error(err)
		}
	})
}
