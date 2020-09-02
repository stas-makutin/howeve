package httpsrv

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/log"
	"github.com/stas-makutin/howeve/tasks"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"golang.org/x/net/netutil"
)

const defaultHTTPPort = 8180

// Task struct
type Task struct {
	cfg      *config.HTTPServerConfig
	listener *net.Listener
	server   *http.Server
	stopWg   sync.WaitGroup
}

// NewTask func
func NewTask() *Task {
	t := &Task{}
	config.AddReader(t.readConfig)
	config.AddWriter(t.writeConfig)
	return t
}

func (t *Task) readConfig(cfg *config.Config, cfgError config.Error) {
	t.cfg = cfg.HTTPServer
	if t.cfg == nil {
		return
	}
	if t.cfg.Port != 0 && (t.cfg.Port < 1 || t.cfg.Port > 65535) {
		cfgError("httpServer.port must be between 1 and 65535.")
	}
}

func (t *Task) writeConfig(cfg *config.Config) {
	cfg.HTTPServer = t.cfg
}

// Open func
func (t *Task) Open(ctx *tasks.ServiceTaskContext) error {

	port := defaultHTTPPort
	var readTimeout, readHeaderTimeout, writeTimeout, idleTimeout uint
	var maxHeaderBytes int
	if t.cfg != nil {
		if t.cfg.Port != 0 {
			port = t.cfg.Port
		}
		readTimeout = t.cfg.ReadTimeout
		readHeaderTimeout = t.cfg.ReadHeaderTimeout
		writeTimeout = t.cfg.WriteTimeout
		idleTimeout = t.cfg.IdleTimeout
		maxHeaderBytes = int(t.cfg.MaxHeaderBytes)
	}

	router := http.NewServeMux()

	setupRoutes(router)

	router.Handle("/socket", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			// handle error
		}
		go func() {
			defer conn.Close()

			for {
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					// handle error
				}
				err = wsutil.WriteServerMessage(conn, op, msg)
				if err != nil {
					// handle error
				}
			}
		}()
	}))

	var handler http.Handler = router
	if log.Enabled() {
		handler = httpLogHandler()(handler)
	}

	server := http.Server{
		Handler:           handler,
		ReadTimeout:       time.Millisecond * time.Duration(readTimeout),
		ReadHeaderTimeout: time.Millisecond * time.Duration(readHeaderTimeout),
		WriteTimeout:      time.Millisecond * time.Duration(writeTimeout),
		IdleTimeout:       time.Millisecond * time.Duration(idleTimeout),
		MaxHeaderBytes:    maxHeaderBytes,
		ErrorLog:          ctx.Log,
	}

	// create TCP listener
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return fmt.Errorf("Listen on %v port failed: %v", port, err)
	}

	// apply concurrent connections limit
	if t.cfg != nil && t.cfg.MaxConnections > 0 {
		listener = netutil.LimitListener(listener, int(t.cfg.MaxConnections))
	}

	t.listener = &listener
	t.server = &server
	t.stopWg.Add(1)
	go func() {
		restart := false
		err = t.server.Serve(*t.listener)
		if err != nil && err != http.ErrServerClosed {
			ctx.Log.Printf("HTTP server failure: %v", err)
			restart = true
		}
		t.stopWg.Done()
		if restart {
			tasks.StopServiceTasks()
		}
	}()

	return nil
}

// Close func
func (t *Task) Close(ctx *tasks.ServiceTaskContext) error {
	return nil
}

// Stop func
func (t *Task) Stop(ctx *tasks.ServiceTaskContext) {
	if t.server != nil {
		err := t.server.Shutdown(context.Background())
		if err != nil {
			ctx.Log.Printf("HTTP server stopping failure: %v", err)
		}
		t.server = nil
	}
	t.stopWg.Wait()
	if t.listener != nil {
		err := (*t.listener).Close()
		if err != nil {
			ctx.Log.Printf("HTTP server listener closing failure: %v", err)
		}
		t.listener = nil
	}
}
