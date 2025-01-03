package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (h Handler) UserAuth(next http.Handler) http.Handler {
	return CustomHandler(func(w http.ResponseWriter, r *http.Request) http.Handler {
		user := r.PathValue("user")
		if user == "" || user == "test" {
			return Error(403, errors.New("Forbidden"))
		}
		return next
	})
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
			fallthrough
		case status < 500:
			h.Logger.Debug(fmt.Sprintf("%s %s", r.Method, r.URL),
				"status", rw.status,
				"duration", time.Since(start))
		default:
			h.Logger.Warn(fmt.Sprintf("%s %s", r.Method, r.URL),
				"status", rw.status,
				"body", rw.writer.String(),
				"duration", time.Since(start))

		}
	})
}

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
