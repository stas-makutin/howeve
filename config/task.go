package config

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/stas-makutin/howeve/tasks"
	"gopkg.in/yaml.v3"
)

// Task struct
type Task struct {
	stopCh     chan struct{}
	watcher    *fsnotify.Watcher
	updateLock uint32
}

// NewTask func
func NewTask() *Task {
	return &Task{
		stopCh: make(chan struct{}),
	}
}

// Open - ServiceTask::Open method
func (t *Task) Open(ctx *tasks.ServiceTaskContext) error {
	if len(ctx.Args) <= 0 {
		return fmt.Errorf("the path to configuration file is not specified")
	}
	cfgFile := ctx.Args[0]
	workingDirectory := ""
	var failure error

	if config, err := readConfig(cfgFile); err == nil {
		ctx.Log.Print("configuration file loaded successfully")

		workingDirectory = config.WorkingDirectory
		if workingDirectory != "" {
			if err := os.Chdir(workingDirectory); err != nil {
				failure = fmt.Errorf("unable to change working directory, reason: %v", err)
			}
		}
	} else {
		failure = err
	}

	if watcher, err := fsnotify.NewWatcher(); err != nil {
		ctx.Log.Printf("unable to create config file watcher, reason: %v", err)
	} else {
		t.watcher = watcher
		ctx.Wg.Add(1)
		go t.watch(&ctx.Wg)
		if err = t.watcher.Add(cfgFile); err != nil {
			t.watcher.Close()
			t.watcher = nil
			ctx.Log.Printf("unable to start config file watcher, reason: %v", err)
		}
	}
	if t.watcher == nil {
		ctx.Wg.Add(1)
		go t.watchFallback(&ctx.Wg, cfgFile)
	}

	writeConfiguration = func(restart bool) bool {
		// block file watcher
		atomic.StoreUint32(&t.updateLock, 1)
		defer atomic.StoreUint32(&t.updateLock, 0)

		var cfg = Config{
			WorkingDirectory: workingDirectory,
		}
		for _, w := range writers {
			w(&cfg)
		}

		if !func() bool {
			file, err := os.OpenFile(cfgFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				ctx.Log.Printf("unable to open configuration file for writing: %v", err)
				return false
			}
			defer file.Close()
			err = yaml.NewEncoder(file).Encode(&cfg)
			if err != nil {
				ctx.Log.Printf("unable to write configuration file: %v", err)
				return false
			}
			return true
		}() {
			return false
		}

		if restart {
			tasks.StopServiceTasks()
		}

		return true
	}

	return failure
}

// Close - ServiceTask::Close method
func (t *Task) Close(ctx *tasks.ServiceTaskContext) error {
	select {
	case <-t.stopCh:
	default:
	}
	return nil
}

// Stop - ServiceTask::Stop method
func (t *Task) Stop(ctx *tasks.ServiceTaskContext) {
	select {
	case t.stopCh <- struct{}{}:
	default:
	}
	if t.watcher != nil {
		t.watcher.Close()
		t.watcher = nil
	}
}

func (t *Task) watch(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		event, ok := <-t.watcher.Events
		if !ok {
			return
		}
		if event.Op&fsnotify.Write == fsnotify.Write {
			if atomic.LoadUint32(&t.updateLock) == 0 {
				tasks.StopServiceTasks()
				return
			}
		}
	}
}

func (t *Task) watchFallback(wg *sync.WaitGroup, cfgFile string) {
	defer wg.Done()
	mt := time.Time{}
	init := true
	for {
		trigger := false
		if fi, err := os.Stat(cfgFile); err == nil {
			if !init && !mt.Equal(fi.ModTime()) {
				trigger = true
			}
			mt = fi.ModTime()
		} else {
			if !init && !mt.Equal(time.Time{}) {
				trigger = true
			}
			mt = time.Time{}
		}
		init = false

		if trigger {
			if atomic.LoadUint32(&t.updateLock) == 0 {
				tasks.StopServiceTasks()
				return
			}
		}

		select {
		case <-t.stopCh:
			return
		case <-time.After(time.Second * 5):
		}
	}
}
