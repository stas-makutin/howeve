package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
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

var writeConfiguration func(restart bool) bool

func addConfigReader(r configReader) {
	configReaders = append(configReaders, r)
}

func addConfigWriter(w configWriter) {
	configWriters = append(configWriters, w)
}

type configTask struct {
	stopCh     chan struct{}
	watcher    *fsnotify.Watcher
	updateLock uint32
}

func newConfigTask() *configTask {
	return &configTask{
		stopCh: make(chan struct{}),
	}
}

func (t *configTask) open(ctx *serviceTaskContext) error {
	if len(ctx.args) <= 0 {
		return fmt.Errorf("the path to configuration file is not specified")
	}
	cfgFile := ctx.args[0]

	workingDirectory := ""
	errMsg := ""
	for {
		config, err := readConfig(cfgFile)
		if err == nil {
			ctx.log.Print("configuration file loaded successfully")
			workingDirectory = config.WorkingDirectory
			break
		}

		msg := err.Error()
		if msg != errMsg {
			errMsg = msg
			ctx.log.Print(errMsg)
		}

		select {
		case <-t.stopCh:
			return err
		case <-time.After(time.Second * 5):
		}
	}

	var err error
	if t.watcher, err = fsnotify.NewWatcher(); err != nil {
		t.watcher = nil
		ctx.log.Printf("unable to create config file watcher, reason: %v", err)
	} else {
		ctx.wg.Add(1)
		go t.watch(&ctx.wg)
		if err = t.watcher.Add(cfgFile); err != nil {
			t.watcher.Close()
			t.watcher = nil
			ctx.log.Printf("unable to start config file watcher, reason: %v", err)
		}
	}

	writeConfiguration = func(restart bool) bool {
		// block file watcher
		atomic.StoreUint32(&t.updateLock, 1)
		defer atomic.StoreUint32(&t.updateLock, 0)

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
			stopServiceTasks()
		}

		return true
	}

	return nil
}

func (t *configTask) close(ctx *serviceTaskContext) error {
	select {
	case <-t.stopCh:
	default:
	}
	return nil
}

func (t *configTask) stop(ctx *serviceTaskContext) {
	select {
	case t.stopCh <- struct{}{}:
	default:
	}
	if t.watcher != nil {
		t.watcher.Close()
		t.watcher = nil
	}
}

func (t *configTask) watch(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case event, ok := <-t.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if atomic.LoadUint32(&t.updateLock) == 0 {
					fmt.Println("watch - event")
					stopServiceTasks()
					return
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

	if config.WorkingDirectory != "" {
		if err := os.Chdir(config.WorkingDirectory); err != nil {
			return nil, fmt.Errorf("unable to change working directory: %v", err)
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
