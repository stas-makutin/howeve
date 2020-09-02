package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/kardianos/service"
	"github.com/stas-makutin/howeve/tasks"
)

const appName = "howeve"

func usage() {
	fmt.Printf(
		`Howeve - Home Automation service.

Usage: [Action] [Options]

Actions:

install [option]
  Install the service
uninstall
  Uninstall the service
start
  Start the service
stop
  Stop the service
restart
  Restart the service
run [option]
  Execute as console/terminal application

Options:
-h, --help
  Print this message
-c, --config <file name>
  Path to configuration yaml file. Works only for install and run actions.
  Default: %v
`,
		defaultConfigFile(),
	)
}

type application struct {
	configFile  string
	logger      service.Logger
	stopService func()
	stopping    bool
	stopped     sync.Mutex
}

// Error logger writer
func (app *application) Write(p []byte) (n int, err error) {
	err = app.logger.Error(string(p))
	if err != nil {
		return
	}
	n = len(p)
	return
}

func (app *application) parseCommandLine(interactive bool) (action string) {
	argc := len(os.Args)
	for i := 1; i < argc; i++ {
		arg := os.Args[i]
		switch arg {
		case "-h", "--help":
			if interactive {
				usage()
				os.Exit(100)
			}
		case "-c", "--config":
			if i++; i < argc {
				app.configFile = os.Args[i]
			}
		default:
			if action == "" {
				action = arg
			}
		}
	}
	return
}

func (app *application) run() {
	app.logger.Info(appName + " started with configuration file " + app.configFile)

	tasks.RunServiceTasks(log.New(app, "", 0), app.configFile)

	if !app.stopping {
		app.stopService()
	}
	app.logger.Info(appName + " stopped")
	app.stopped.Unlock()
}

func (app *application) Start(s service.Service) error {
	if app.configFile == "" {
		app.parseCommandLine(false)
		if app.configFile == "" {
			app.configFile = defaultConfigFile()
		}
	}
	go app.run()
	return nil
}

func (app *application) Stop(s service.Service) error {
	app.stopping = true
	tasks.EndServiceTasks()
	return nil
}

func main() {
	app := &application{}
	action := app.parseCommandLine(service.Interactive())

	var arguments []string
	if app.configFile != "" {
		arguments = []string{"--config", app.configFile}
	}

	svcConfig := &service.Config{
		Name:        appName,
		DisplayName: "Home Automation Service",
		Description: "Home Automation Service.",
		Arguments:   arguments,
	}

	svc, err := service.New(app, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	app.logger, err = svc.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	app.stopped.Lock()
	if service.Interactive() {
		if action == "run" {
			app.stopService = func() {}
			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc, os.Interrupt)
			go func() {
				<-sigc
				if err := app.Stop(svc); err != nil {
					app.logger.Error(err)
					os.Exit(1)
				}
			}()
			err = app.Start(svc)
			if err != nil {
				app.logger.Error(err)
				os.Exit(1)
			}
			app.stopped.Lock()
			os.Exit(0)
		} else if action != "" {
			err = service.Control(svc, action)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Please specify the Action. Run with --help for details.")
		}
		return
	}
	app.stopService = func() {
		svc.Stop()
	}
	err = svc.Run()
	if err != nil {
		app.logger.Error(err)
	}
}
