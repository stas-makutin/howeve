package main

import (
	"github.com/stas-makutin/howeve/howeve/config"
	"github.com/stas-makutin/howeve/howeve/httpsrv"
	"github.com/stas-makutin/howeve/howeve/log"
	"github.com/stas-makutin/howeve/howeve/tasks"
)

func init() {
	tasks.ServiceTasks = []tasks.ServiceTaskEntry{
		{Name: "Configuration", Task: config.NewTask()},
		{Name: "Log", Task: log.NewTask(appName)},
		{Name: "HTTP server", Task: httpsrv.NewTask()},
	}
}
