// +build !windows

package main

import "path/filepath"

// NewLine is platform-specific new line character sequence
const NewLine = "\n"

func defaultConfigFile() string {
	return filepath.Join("/etc", appName, appName+".yml")
}
