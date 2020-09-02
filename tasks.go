package main

import (
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/httpsrv"
	"github.com/stas-makutin/howeve/log"
	"github.com/stas-makutin/howeve/tasks"
)

func init() {
	tasks.ServiceTasks = []tasks.ServiceTaskEntry{
		{Name: "Configuration", Task: config.NewTask()},
		{Name: "Log", Task: log.NewTask(appName)},
		{Name: "HTTP server", Task: httpsrv.NewTask()},
	}
}
