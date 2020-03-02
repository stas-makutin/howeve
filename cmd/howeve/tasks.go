package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

type serviceTaskContext struct {
	args []string
	log  *log.Logger
	wg   sync.WaitGroup
}

type serviceTask interface {
	start(ctx *serviceTaskContext) error
	stop(ctx *serviceTaskContext) error
}

type serviceTaskEntry struct {
	name string
	task serviceTask
}

var serviceTasks []serviceTaskEntry
var serviceTaskCtx serviceTaskContext
var syncTaskOp sync.Mutex
var taskStopping bool = false

func runServiceTasks(errorLog *log.Logger, cfgFile string) error {
	// initialize the context
	serviceTaskCtx.log = errorLog
	serviceTaskCtx.args = []string{cfgFile}

	// run tasks
	run := true
	for run {
		var failure error = nil

		syncTaskOp.Lock()
		run = func() bool {
			defer syncTaskOp.Unlock()

			if taskStopping {
				return false
			}

			var errStr strings.Builder
			for i, te := range serviceTasks {
				failure = te.task.start(&serviceTaskCtx)
				if failure != nil {
					writeStringln(&errStr, fmt.Sprintf("%v task failed to start: %v", te.name, failure))
					for j := i - 1; j >= 0; j-- {
						ste := serviceTasks[j]
						if es := ste.task.stop(&serviceTaskCtx); es != nil {
							writeStringln(&errStr, fmt.Sprintf("%v task failed to stop: %v", ste.name, es))
						}
					}
				}
			}
			if errStr.Len() > 0 {
				serviceTaskCtx.log.Print(errStr.String())
			}

			return true
		}()

		serviceTaskCtx.wg.Wait()

		if failure != nil {
			return failure
		}
	}

	return nil
}

func stopServiceTasks(restart bool) error {
	syncTaskOp.Lock()
	return func() error {
		defer syncTaskOp.Unlock()

		var failure error = nil

		taskStopping = !restart

		var errStr strings.Builder
		for j := len(serviceTasks) - 1; j >= 0; j-- {
			te := serviceTasks[j]
			if err := te.task.stop(&serviceTaskCtx); err != nil {
				if failure == nil {
					failure = err
				}
				writeStringln(&errStr, fmt.Sprintf("%v task failed to stop: %v", te.name, err))
			}
		}
		if errStr.Len() > 0 {
			serviceTaskCtx.log.Print(errStr.String())
		}

		return failure
	}()
}
