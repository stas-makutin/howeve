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
	serviceTaskEntry{"Log", newLogTask()},
}

var serviceTaskCtx serviceTaskContext
var serviceTaskStop uint32
var serviceTaskClose uint32

func runServiceTasks(errorLog *log.Logger, cfgFile string) {

	// initialize the context
	serviceTaskCtx.log = errorLog
	serviceTaskCtx.args = []string{cfgFile}

	// run tasks
	prevErrorMsg := ""
	for atomic.LoadUint32(&serviceTaskClose) == 0 {
		var index int = -1

		errorMsg := ""
		for i, te := range serviceTasks {
			for atomic.LoadUint32(&serviceTaskStop) != 0 {
				break
			}
			if err := te.task.open(&serviceTaskCtx); err != nil {
				errorMsg = fmt.Sprintf("%v task failed to open: %v", te.name, err)
				break
			}
			index = i
		}
		if errorMsg != "" {
			if prevErrorMsg != errorMsg {
				prevErrorMsg = errorMsg
				serviceTaskCtx.log.Print(errorMsg)
			}
		} else {
			prevErrorMsg = ""
		}

		serviceTaskCtx.wg.Wait()

		atomic.StoreUint32(&serviceTaskStop, 0)

		var closeErrorMsg strings.Builder
		for ; index >= 0; index-- {
			te := serviceTasks[index]
			if err := te.task.close(&serviceTaskCtx); err != nil {
				writeStringln(&closeErrorMsg, fmt.Sprintf("%v task failed to close: %v", te.name, err))
			}
		}
		if closeErrorMsg.Len() > 0 {
			serviceTaskCtx.log.Print(closeErrorMsg.String())
			break // fatal error, stopping
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
