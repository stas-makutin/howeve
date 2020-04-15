// +build windows

package main

import (
	"os"
	"path/filepath"
)

// NewLine is platform-specific new line character sequence
const NewLine = "\r\n"

func defaultConfigFile() string {
	return filepath.Join(filepath.Dir(os.Args[0]), appName+".yml")
}
