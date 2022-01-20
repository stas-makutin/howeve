package config

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
	Dir        string       `yaml:"dir,omitempty" json:"dir,omitempty"`
	File       string       `yaml:"file,omitempty" json:"file,omitempty"`
	DirMode    FileMode     `yaml:"dirMode,omitempty" json:"dirMode,omitempty"`
	FileMode   FileMode     `yaml:"fileMode,omitempty" json:"fileMode,omitempty"`
	MaxSize    SizeType     `yaml:"maxSize,omitempty" json:"maxSize,omitempty"`
	MaxAge     DurationType `yaml:"maxAge,omitempty" json:"maxAge,omitempty"` // seconds
	Backups    uint32       `yaml:"backups,omitempty" json:"backups,omitempty"`
	BackupDays uint32       `yaml:"backupDays,omitempty" json:"backupDays,omitempty"`
	Archive    string       `yaml:"archive,omitempty" json:"archive,omitempty"`
}

// HTTPServerConfig defines configuration entries for the HTTP server
type HTTPServerConfig struct {
	Port              int          `yaml:"port,omitempty" json:"port,omitempty"`
	MaxConnections    uint         `yaml:"maxConnections,omitempty" json:"maxConnections,omitempty"`
	ReadTimeout       DurationType `yaml:"readTimeout,omitempty" json:"readTimeout,omitempty"`
	ReadHeaderTimeout DurationType `yaml:"readHeaderTimeout,omitempty" json:"readHeaderTimeout,omitempty"`
	WriteTimeout      DurationType `yaml:"writeTimeout,omitempty" json:"writeTimeout,omitempty"`
	IdleTimeout       DurationType `yaml:"idleTimeout,omitempty" json:"idleTimeout,omitempty"`
	MaxHeaderBytes    SizeType     `yaml:"maxHeaderBytes,omitempty" json:"maxHeaderBytes,omitempty"`
	Assets            []HTTPAsset  `yaml:"assets,omitempty" json:"assets,omitempty"`
}

type HTTPAsset struct {
	Route string `yaml:"route,omitempty" json:"route,omitempty"`
	Path  string `yaml:"path,omitempty" json:"path,omitempty"`
}

// ServiceConfig defines configuration of active services
type ServiceConfig struct {
	Alias     string            `yaml:"alias,omitempty" json:"alias,omitempty"`
	Protocol  string            `yaml:"protocol" json:"protocol"`
	Transport string            `yaml:"transport" json:"transport"`
	Entry     string            `yaml:"entry,omitempty" json:"entry,omitempty"`
	Params    map[string]string `yaml:"params,omitempty" json:"params,omitempty"`
}
