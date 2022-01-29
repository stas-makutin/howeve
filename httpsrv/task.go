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
	cancel   context.CancelFunc
	hc       handlerContext
}

// NewTask func
func NewTask() *Task {
	t := &Task{}
	t.hc.stopCh = make(chan struct{})
	close(t.hc.stopCh)
	config.AddReader(t.readConfig)
	config.AddWriter(t.writeConfig)
	return t
}

func (t *Task) readConfig(cfg *config.Config, cfgError config.Error) {
	t.hc.cfg = cfg.HTTPServer
	if t.hc.cfg == nil {
		return
	}
	if t.hc.cfg.Port != 0 && (t.hc.cfg.Port < 1 || t.hc.cfg.Port > 65535) {
		cfgError("httpServer.port must be between 1 and 65535")
	}
}

func (t *Task) writeConfig(cfg *config.Config) {
	cfg.HTTPServer = t.hc.cfg
}

// Open func
func (t *Task) Open(ctx *tasks.ServiceTaskContext) error {

	port := defaultHTTPPort
	var readTimeout, readHeaderTimeout, writeTimeout, idleTimeout time.Duration
	var maxHeaderBytes int
	if t.hc.cfg != nil {
		if t.hc.cfg.Port != 0 {
			port = t.hc.cfg.Port
		}
		readTimeout = t.hc.cfg.ReadTimeout.Value()
		readHeaderTimeout = t.hc.cfg.ReadHeaderTimeout.Value()
		writeTimeout = t.hc.cfg.WriteTimeout.Value()
		idleTimeout = t.hc.cfg.IdleTimeout.Value()
		maxHeaderBytes = int(t.hc.cfg.MaxHeaderBytes.Value())
	}

	router := http.NewServeMux()

	setupRoutes(router, t.hc.cfg.Assets)

	var handler http.Handler = router
	if log.Enabled() {
		handler = logHandler(handler)
	}

	baseCtx, cancel := context.WithCancel(context.WithValue(context.Background(), handlerContextKey, &t.hc))
	t.cancel = cancel

	server := http.Server{
		Handler:           handler,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
		ErrorLog:          ctx.Log,
		BaseContext:       func(listener net.Listener) context.Context { return baseCtx },
	}

	// create TCP listener
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return fmt.Errorf("listen on %v port failed: %v", port, err)
	}

	// apply concurrent connections limit
	if t.hc.cfg != nil && t.hc.cfg.MaxConnections > 0 {
		listener = netutil.LimitListener(listener, int(t.hc.cfg.MaxConnections))
	}

	t.listener = &listener
	t.server = &server
	ctx.Wg.Add(1)
	t.hc.stopCh = make(chan struct{})
	go func() {
		restart := false
		err = t.server.Serve(*t.listener)
		if err != nil && err != http.ErrServerClosed {
			ctx.Log.Printf("HTTP server failure: %v", err)
			restart = true
		}
		close(t.hc.stopCh)
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
		t.cancel()
		err := t.server.Shutdown(context.Background())
		if err != nil {
			ctx.Log.Printf("HTTP server stopping failure: %v", err)
		}
		t.cancel = nil
		t.server = nil
	}

	<-t.hc.stopCh

	t.hc.handlerWg.Wait()
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
