package httpsrv

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/stas-makutin/howeve/log"
)

type logResponseWriter struct {
	http.ResponseWriter
	statusCode    int
	contentLength int64
}

func (w *logResponseWriter) WriteHeader(status int) {
	w.statusCode = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *logResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.contentLength += int64(n)
	return n, err
}

func (w *logResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, fmt.Errorf("the instance of http.ResponseWriter is not http.Hijacker")
}

type httpContextKey int

const httpLogFieldsKey httpContextKey = 0

type httpLogFields struct {
	fields []string
}

func appendLogFields(r *http.Request, vals ...string) {
	if fields, ok := r.Context().Value(httpLogFieldsKey).(*httpLogFields); ok {
		fields.fields = append(fields.fields, vals...)
	}
}

func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().Local()
		lrw := &logResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		fields := &httpLogFields{}
		ctx := context.WithValue(r.Context(), httpLogFieldsKey, fields)
		defer func() {
			log.Report(append([]string{
				log.SrcHTTP,
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
