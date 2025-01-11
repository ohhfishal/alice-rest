package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (h Handler) UserAuth(next http.Handler) http.Handler {
	// TODO: Implement
	return CustomHandler(func(w http.ResponseWriter, r *http.Request) http.Handler {
		user := r.PathValue("user")
		if user == "" || user == "test" {
			return Error(403, errors.New("Forbidden"))
		}
		return next
	})
}

func (h Handler) SetupContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := context.WithValue(r.Context(), "start", start)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h Handler) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := responseWriter{ResponseWriter: w}
		next.ServeHTTP(&rw, r)

		status := rw.Status()

		ctx := context.WithValue(r.Context(), "status", status)
		ctxErr := ctx.Err()

		switch {
		case ctxErr != nil:
			h.Logger.ErrorContext(ctx, ctxErr.Error())
		case status < 300:
			fallthrough
		case status < 400:
			fallthrough
		case status < 500:
			h.Logger.DebugContext(ctx, fmt.Sprintf("%s %s", r.Method, r.URL))
		default:
			h.Logger.WarnContext(ctx, fmt.Sprintf("%s %s", r.Method, r.URL), "body", rw.writer.String())
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
