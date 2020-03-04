package main

import (
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type configError func(msg string)

type configReader func(cfg *Config, cfgError configError)

type configWriter func(cfg *Config)

var configReaders []configReader
var configWriters []configWriter

var writeConfigurationLock int32
var writeConfiguration func(restart bool) bool

func addConfigReader(r configReader) {
	configReaders = append(configReaders, r)
}

func addConfigWriter(w configWriter) {
	configWriters = append(configWriters, w)
}

type configTask struct {
	watcher *fsnotify.Watcher
}

func (t *configTask) start(ctx *serviceTaskContext) error {
	if len(ctx.args) <= 0 {
		return fmt.Errorf("the path to configuration file is not specified")
	}
	cfgFile := ctx.args[0]

	workingDirectory := ""
	for {
		config, err := readConfig(cfgFile)
		if err == nil {
			workingDirectory = config.WorkingDirectory
			break
		}
		time.Sleep(time.Second * 5)
	}

	var err error
	if t.watcher, err = fsnotify.NewWatcher(); err != nil {
		t.watcher = nil
		ctx.log.Printf("unable to create config file watcher, reason: %v", err)
	} else {
		go t.watch()
		if err = t.watcher.Add(cfgFile); err != nil {
			t.watcher.Close()
			t.watcher = nil
			ctx.log.Printf("unable to start config file watcher, reason: %v", err)
		} else {
			ctx.wg.Add(1)
		}
	}

	writeConfiguration = func(restart bool) bool {
		// block file watcher
		atomic.CompareAndSwapInt32(&writeConfigurationLock, 0, 1)

		var cfg = Config{
			WorkingDirectory: workingDirectory,
		}
		for _, w := range configWriters {
			w(&cfg)
		}

		if !func() bool {
			file, err := os.OpenFile(cfgFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				ctx.log.Printf("unable to open configuration file for writing: %v", err)
				return false
			}
			defer file.Close()
			err = yaml.NewEncoder(file).Encode(&cfg)
			if err != nil {
				ctx.log.Printf("unable to write configuration file: %v", err)
				return false
			}
			return true
		}() {
			return false
		}

		if restart {
			go func() {
				stopServiceTasks(true)
			}()
		}

		return true
	}

	return nil
}

func (t *configTask) stop(ctx *serviceTaskContext) error {
	if t.watcher != nil {
		t.watcher.Close()
		t.watcher = nil
		ctx.wg.Done()
	}
	return nil
}

func (t *configTask) watch() {
	for {
		select {
		case event, ok := <-t.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if !atomic.CompareAndSwapInt32(&writeConfigurationLock, 1, 0) {
					stopServiceTasks(true)
				}
			}
		}
	}
}

func readConfig(cfgFile string) (*Config, error) {
	var config Config

	err := func() error {
		file, err := os.Open(cfgFile)
		if err != nil {
			return fmt.Errorf("unable to open configuration file '%v': %v", cfgFile, err)
		}
		defer file.Close()
		err = yaml.NewDecoder(file).Decode(&config)
		if err != nil {
			return fmt.Errorf("unable to parse configuration file '%v': %v", cfgFile, err)
		}
		return nil
	}()
	if err != nil {
		return nil, err
	}

	if config.WorkingDirectory != "" {
		if err := os.Chdir(config.WorkingDirectory); err != nil {
			return nil, fmt.Errorf("unable to change working directory to '%v': %v", config.WorkingDirectory, err)
		}
	}

	var errStr strings.Builder
	ce := func(msg string) {
		errStr.WriteString(NewLine + msg)
	}

	for _, r := range configReaders {
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
