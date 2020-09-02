package tasks

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/stas-makutin/howeve/utils"
)

// ServiceTaskContext struct
type ServiceTaskContext struct {
	Args []string
	Log  *log.Logger
	Wg   sync.WaitGroup
}

// ServiceTask interface
type ServiceTask interface {
	Open(ctx *ServiceTaskContext) error
	Close(ctx *ServiceTaskContext) error
	Stop(ctx *ServiceTaskContext)
}

// ServiceTaskEntry struct
type ServiceTaskEntry struct {
	Name string
	Task ServiceTask
}

// ServiceTasks - slice of known tasks to run
var ServiceTasks []ServiceTaskEntry
var serviceTaskCtx ServiceTaskContext
var serviceTaskStop uint32
var serviceTaskClose uint32

// RunServiceTasks func
func RunServiceTasks(errorLog *log.Logger, cfgFile string) {

	// initialize the context
	serviceTaskCtx.Log = errorLog
	serviceTaskCtx.Args = []string{cfgFile}

	// run tasks
	prevErrorMsg := ""
	for atomic.LoadUint32(&serviceTaskClose) == 0 {
		var index int = -1

		errorMsg := ""
		for i, te := range ServiceTasks {
			for atomic.LoadUint32(&serviceTaskStop) != 0 {
				break
			}
			if err := te.Task.Open(&serviceTaskCtx); err != nil {
				errorMsg = fmt.Sprintf("%v task failed to open: %v", te.Name, err)
				break
			}
			index = i
		}
		if errorMsg != "" {
			if prevErrorMsg != errorMsg {
				prevErrorMsg = errorMsg
				serviceTaskCtx.Log.Print(errorMsg)
			}
		} else {
			prevErrorMsg = ""
		}

		serviceTaskCtx.Wg.Wait()

		atomic.StoreUint32(&serviceTaskStop, 0)

		var closeErrorMsg strings.Builder
		for ; index >= 0; index-- {
			te := ServiceTasks[index]
			if err := te.Task.Close(&serviceTaskCtx); err != nil {
				utils.WriteStringln(&closeErrorMsg, fmt.Sprintf("%v task failed to close: %v", te.Name, err))
			}
		}
		if closeErrorMsg.Len() > 0 {
			serviceTaskCtx.Log.Print(closeErrorMsg.String())
			break // fatal error, stopping
		}
	}
}

// StopServiceTasks func
func StopServiceTasks() {
	atomic.StoreUint32(&serviceTaskStop, 1)
	for _, te := range ServiceTasks {
		te.Task.Stop(&serviceTaskCtx)
	}
}

// EndServiceTasks func
func EndServiceTasks() {
	atomic.StoreUint32(&serviceTaskClose, 1)
	StopServiceTasks()
}
