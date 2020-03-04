package main

// Config is a collection of top-level entries of configuration file
type Config struct {
	WorkingDirectory string `yaml:"workingDir,omitempty"`
}
