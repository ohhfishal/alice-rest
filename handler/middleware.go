package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	writer      strings.Builder
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) Write(bytes []byte) (int, error) {
	rw.writer.Write(bytes)
	return rw.ResponseWriter.Write(bytes)
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

func (h Handler) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := responseWriter{ResponseWriter: w}
		start := time.Now()
		next.ServeHTTP(&rw, r)
		status := rw.status
		switch {
		case status < 300:
			fallthrough
		case status < 400:
			h.Logger.Debug(fmt.Sprintf("%s %s", r.Method, r.URL),
				"status", rw.status,
				"duration", time.Since(start))
		case status < 500:
			fallthrough
		default:
			h.Logger.Warn(fmt.Sprintf("%s %s", r.Method, r.URL),
				"status", rw.status,
				"body", rw.writer.String(),
				"duration", time.Since(start))

		}
	})
}
