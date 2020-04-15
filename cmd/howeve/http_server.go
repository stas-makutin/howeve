package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/netutil"
)

const defaultHTTPPort = 8180

type httpServerTask struct {
	cfg      *HTTPServerConfig
	listener *net.Listener
	server   *http.Server
	stopWg   sync.WaitGroup
}

func newHTTPServerTask() *httpServerTask {
	t := &httpServerTask{}
	addConfigReader(t.readConfig)
	addConfigWriter(t.writeConfig)
	return t
}

func (t *httpServerTask) readConfig(cfg *Config, cfgError configError) {
	t.cfg = cfg.HTTPServer
	if t.cfg == nil {
		return
	}
	if t.cfg.Port != 0 && (t.cfg.Port < 1 || t.cfg.Port > 65535) {
		cfgError("httpServer.port must be between 1 and 65535.")
	}
}

func (t *httpServerTask) writeConfig(cfg *Config) {
	cfg.HTTPServer = t.cfg
}

func (t *httpServerTask) open(ctx *serviceTaskContext) error {

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

	var handler http.Handler = router
	/*
		if config.HttpServer.Log != nil {
			handler = logHandler(errorLog)(handler)
		}
	*/

	server := http.Server{
		Handler:           handler,
		ReadTimeout:       time.Millisecond * time.Duration(readTimeout),
		ReadHeaderTimeout: time.Millisecond * time.Duration(readHeaderTimeout),
		WriteTimeout:      time.Millisecond * time.Duration(writeTimeout),
		IdleTimeout:       time.Millisecond * time.Duration(idleTimeout),
		MaxHeaderBytes:    maxHeaderBytes,
		ErrorLog:          ctx.log,
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
			ctx.log.Printf("HTTP server failure: %v", err)
			restart = true
		}
		t.stopWg.Done()
		if restart {
			stopServiceTasks()
		}
	}()

	return nil
}

func (t *httpServerTask) close(ctx *serviceTaskContext) error {
	return nil
}

func (t *httpServerTask) stop(ctx *serviceTaskContext) {
	if t.server != nil {
		err := t.server.Shutdown(context.Background())
		if err != nil {
			ctx.log.Printf("HTTP server stopping failure: %v", err)
		}
		t.server = nil
	}
	t.stopWg.Wait()
	if t.listener != nil {
		err := (*t.listener).Close()
		if err != nil {
			ctx.log.Printf("HTTP server listener closing failure: %v", err)
		}
		t.listener = nil
	}
}
