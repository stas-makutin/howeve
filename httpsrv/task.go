package httpsrv

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/log"
	"github.com/stas-makutin/howeve/tasks"

	"golang.org/x/net/netutil"
)

const defaultHTTPPort = 8180

// Task struct
type Task struct {
	listener *net.Listener
	server   *http.Server
	hctx     handlerContext
}

// NewTask func
func NewTask() *Task {
	t := &Task{}
	config.AddReader(t.readConfig)
	config.AddWriter(t.writeConfig)
	return t
}

func (t *Task) readConfig(cfg *config.Config, cfgError config.Error) {
	t.hctx.cfg = cfg.HTTPServer
	if t.hctx.cfg == nil {
		return
	}
	if t.hctx.cfg.Port != 0 && (t.hctx.cfg.Port < 1 || t.hctx.cfg.Port > 65535) {
		cfgError("httpServer.port must be between 1 and 65535.")
	}
}

func (t *Task) writeConfig(cfg *config.Config) {
	cfg.HTTPServer = t.hctx.cfg
}

// Open func
func (t *Task) Open(ctx *tasks.ServiceTaskContext) error {

	port := defaultHTTPPort
	var readTimeout, readHeaderTimeout, writeTimeout, idleTimeout uint
	var maxHeaderBytes int
	if t.hctx.cfg != nil {
		if t.hctx.cfg.Port != 0 {
			port = t.hctx.cfg.Port
		}
		readTimeout = t.hctx.cfg.ReadTimeout
		readHeaderTimeout = t.hctx.cfg.ReadHeaderTimeout
		writeTimeout = t.hctx.cfg.WriteTimeout
		idleTimeout = t.hctx.cfg.IdleTimeout
		maxHeaderBytes = int(t.hctx.cfg.MaxHeaderBytes)
	}

	router := http.NewServeMux()

	setupRoutes(router)

	var handler http.Handler = router
	if log.Enabled() {
		handler = httpLogHandler()(handler)
	}

	baseCtx := context.WithValue(context.Background(), handlerContextKey, &t.hctx)

	server := http.Server{
		Handler:           handler,
		ReadTimeout:       time.Millisecond * time.Duration(readTimeout),
		ReadHeaderTimeout: time.Millisecond * time.Duration(readHeaderTimeout),
		WriteTimeout:      time.Millisecond * time.Duration(writeTimeout),
		IdleTimeout:       time.Millisecond * time.Duration(idleTimeout),
		MaxHeaderBytes:    maxHeaderBytes,
		ErrorLog:          ctx.Log,
		BaseContext:       func(listener net.Listener) context.Context { return baseCtx },
	}

	// create TCP listener
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return fmt.Errorf("Listen on %v port failed: %v", port, err)
	}

	// apply concurrent connections limit
	if t.hctx.cfg != nil && t.hctx.cfg.MaxConnections > 0 {
		listener = netutil.LimitListener(listener, int(t.hctx.cfg.MaxConnections))
	}

	t.listener = &listener
	t.server = &server
	ctx.Wg.Add(1)
	t.hctx.stopWg.Add(1)
	go func() {
		restart := false
		err = t.server.Serve(*t.listener)
		if err != nil && err != http.ErrServerClosed {
			ctx.Log.Printf("HTTP server failure: %v", err)
			restart = true
		}
		t.hctx.stopWg.Done()
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
	defer ctx.Wg.Done()

	if t.server != nil {
		err := t.server.Shutdown(context.Background())
		if err != nil {
			ctx.Log.Printf("HTTP server stopping failure: %v", err)
		}
		t.server = nil
	}
	t.hctx.stopWg.Wait()
	t.hctx.handlerWg.Wait()
	if t.listener != nil {
		err := (*t.listener).Close()
		if err != nil {
			if operr, ok := err.(*net.OpError); !ok || operr.Op != "close" {
				ctx.Log.Printf("HTTP server listener closing failure: %v", err)
			}
		}
		t.listener = nil
	}
}
