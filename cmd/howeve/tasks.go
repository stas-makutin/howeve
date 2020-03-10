package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
)

type serviceTaskContext struct {
	args []string
	log  *log.Logger
	wg   sync.WaitGroup
}

type serviceTask interface {
	open(ctx *serviceTaskContext) error
	close(ctx *serviceTaskContext) error
	stop(ctx *serviceTaskContext)
}

type serviceTaskEntry struct {
	name string
	task serviceTask
}

var serviceTasks = []serviceTaskEntry{
	serviceTaskEntry{"Configuration", newConfigTask()},
}

var serviceTaskCtx serviceTaskContext
var serviceTaskStop uint32
var serviceTaskClose uint32

func runServiceTasks(errorLog *log.Logger, cfgFile string) {

	// initialize the context
	serviceTaskCtx.log = errorLog
	serviceTaskCtx.args = []string{cfgFile}

	// run tasks
	for atomic.LoadUint32(&serviceTaskClose) == 0 {
		var emsg strings.Builder
		var index int = -1
		var wait bool = true

		for i, te := range serviceTasks {
			for atomic.LoadUint32(&serviceTaskStop) != 0 {
				wait = false
				break
			}
			if err := te.task.open(&serviceTaskCtx); err != nil {
				writeStringln(&emsg, fmt.Sprintf("%v task failed to open: %v", te.name, err))
				wait = false
				break
			}
			index = i
		}

		if wait {
			serviceTaskCtx.wg.Wait()
		}

		atomic.StoreUint32(&serviceTaskStop, 0)

		for ; index >= 0; index-- {
			te := serviceTasks[index]
			if err := te.task.close(&serviceTaskCtx); err != nil {
				writeStringln(&emsg, fmt.Sprintf("%v task failed to close: %v", te.name, err))
			}
		}

		if emsg.Len() > 0 {
			serviceTaskCtx.log.Print(emsg.String())
			return
		}
	}
}

func stopServiceTasks() {
	atomic.StoreUint32(&serviceTaskStop, 1)
	for _, te := range serviceTasks {
		te.task.stop(&serviceTaskCtx)
	}
}

func endServiceTasks() {
	atomic.StoreUint32(&serviceTaskClose, 1)
	stopServiceTasks()
}
