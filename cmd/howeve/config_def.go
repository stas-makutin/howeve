package main

import "os"

// Config is a collection of top-level entries of configuration file
type Config struct {
	WorkingDirectory string     `yaml:"workingDir,omitempty"`
	Log              *LogConfig `yaml:"log,omitempty"`
}

// LogConfig defines configuration entries for the serivce logging
type LogConfig struct {
	Dir        string      `yaml:"dir,omitempty"`
	File       string      `yaml:"file,omitempty"`
	DirMode    os.FileMode `yaml:"dirMode,omitempty"`
	FileMode   os.FileMode `yaml:"fileMode,omitempty"`
	MaxSize    string      `yaml:"maxSize,omitempty"`
	MaxAge     string      `yaml:"maxAge,omitempty"` // seconds
	Backups    uint32      `yaml:"backups,omitempty"`
	BackupDays uint32      `yaml:"backupDays,omitempty"`
	Archive    string      `yaml:"archive,omitempty"`
}
