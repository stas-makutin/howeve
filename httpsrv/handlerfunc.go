package httpsrv

import (
	"net/http"
	"sync"

	"github.com/stas-makutin/howeve/defs"
)

// HandlerContextKeyType handler context field type
type handlerContextKeyType int

// HandlerContextKey - handler context field in the http request context
const handlerContextKey handlerContextKeyType = 0

// HandlerContext struct
type handlerContext struct {
	cfg       *defs.HTTPServerConfig
	handlerWg sync.WaitGroup
	stopCh    chan struct{}
}

type handlerCtxFunc func(http.ResponseWriter, *http.Request, *handlerContext)

func (f handlerCtxFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if hc, ok := r.Context().Value(handlerContextKey).(*handlerContext); ok {
		f(w, r, hc)
	}
}

type handlerFunc func(http.ResponseWriter, *http.Request)

func (f handlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if hc, ok := r.Context().Value(handlerContextKey).(*handlerContext); ok {
		hc.handlerWg.Add(1)
		defer hc.handlerWg.Done()
		f(w, r)
	}
}
