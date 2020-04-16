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

	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Test")
	}))

	var handler http.Handler = router
	if logEnabled() {
		handler = httpLogHandler()(handler)
	}

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

type httpLogResponseWriter struct {
	http.ResponseWriter
	statusCode    int
	contentLength int64
}

func (w *httpLogResponseWriter) WriteHeader(status int) {
	w.statusCode = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *httpLogResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.contentLength += int64(n)
	return n, err
}

type httpContextKey int

const httpLogFieldsKey httpContextKey = 0

type httpLogFields struct {
	fields []string
}

func httpAppendLogField(r *http.Request, val string) {
	if fields, ok := r.Context().Value(httpLogFieldsKey).(httpLogFields); ok {
		fields.fields = append(fields.fields, val)
	}
}

func httpLogHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now().Local()
			lrw := &httpLogResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			fields := httpLogFields{}
			ctx := context.WithValue(r.Context(), httpLogFieldsKey, fields)
			defer func() {
				logr(append([]string{
					logSourceHTTP,
					strconv.FormatInt(int64(time.Now().Local().Sub(start)/time.Millisecond), 10),
					r.RemoteAddr,
					r.Host,
					r.Proto,
					r.Method,
					r.RequestURI,
					strconv.FormatInt(r.ContentLength, 10),
					r.Header.Get("X-Request-Id"),
					strconv.Itoa(lrw.statusCode),
					strconv.FormatInt(lrw.contentLength, 10),
				}, fields.fields...)...)
			}()
			next.ServeHTTP(lrw, r.WithContext(ctx))
		})
	}
}
