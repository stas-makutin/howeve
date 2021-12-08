package config

import "os"

// Config is a collection of top-level entries of configuration file
type Config struct {
	WorkingDirectory string            `yaml:"workingDir,omitempty" json:"workingDir,omitempty"`
	Log              *LogConfig        `yaml:"log,omitempty" json:"log,omitempty"`
	HTTPServer       *HTTPServerConfig `yaml:"httpServer,omitempty" json:"httpServer,omitempty"`
	MessageLog       *MessageLogConfig `yaml:"messageLog,omitempty" json:"messageLog,omitempty"`
	Services         []ServiceConfig   `yaml:"services,omitempty" json:"services,omitempty"`
}

// LogConfig defines configuration entries for the serivce logging
type LogConfig struct {
	Dir        string      `yaml:"dir,omitempty" json:"dir,omitempty"`
	File       string      `yaml:"file,omitempty" json:"file,omitempty"`
	DirMode    os.FileMode `yaml:"dirMode,omitempty" json:"dirMode,omitempty"`
	FileMode   os.FileMode `yaml:"fileMode,omitempty" json:"fileMode,omitempty"`
	MaxSize    string      `yaml:"maxSize,omitempty" json:"maxSize,omitempty"`
	MaxAge     string      `yaml:"maxAge,omitempty" json:"maxAge,omitempty"` // seconds
	Backups    uint32      `yaml:"backups,omitempty" json:"backups,omitempty"`
	BackupDays uint32      `yaml:"backupDays,omitempty" json:"backupDays,omitempty"`
	Archive    string      `yaml:"archive,omitempty" json:"archive,omitempty"`
}

// HTTPServerConfig defines configuration entries for the HTTP server
type HTTPServerConfig struct {
	Port              int    `yaml:"port,omitempty" json:"port,omitempty"`
	MaxConnections    uint   `yaml:"maxConnections,omitempty" json:"maxConnections,omitempty"`
	ReadTimeout       uint   `yaml:"readTimeout,omitempty" json:"readTimeout,omitempty"`             // milliseconds
	ReadHeaderTimeout uint   `yaml:"readHeaderTimeout,omitempty" json:"readHeaderTimeout,omitempty"` // milliseconds
	WriteTimeout      uint   `yaml:"writeTimeout,omitempty" json:"writeTimeout,omitempty"`           // milliseconds
	IdleTimeout       uint   `yaml:"idleTimeout,omitempty" json:"idleTimeout,omitempty"`             // milliseconds
	MaxHeaderBytes    uint32 `yaml:"maxHeaderBytes,omitempty" json:"maxHeaderBytes,omitempty"`
}

// ServiceConfig defines configuration of active services
type ServiceConfig struct {
	Alias     string            `yaml:"alias,omitempty" json:"alias,omitempty"`
	Protocol  string            `yaml:"protocol" json:"protocol"`
	Transport string            `yaml:"transport" json:"transport"`
	Entry     string            `yaml:"entry,omitempty" json:"entry,omitempty"`
	Params    map[string]string `yaml:"params,omitempty" json:"params,omitempty"`
}
