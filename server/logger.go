package server

import (
	"context"
	"io"
	"log/slog"
	"time"
)

func NewLogger(stdout io.Writer, level slog.Level) *slog.Logger {
	base := slog.New(slog.NewJSONHandler(stdout, &slog.HandlerOptions{
		Level: level,
	}))

	handler := LogHandler{
		Handler: base.Handler(),
	}
	return slog.New(handler)
}

type LogHandler struct {
	slog.Handler
}

func (handler LogHandler) Handle(ctx context.Context, r slog.Record) error {
	if start, ok := ctx.Value("start").(time.Time); ok {
		r.AddAttrs(slog.Time("start", start))
		r.AddAttrs(slog.Duration("duration", time.Since(start)))
	}

	if status, ok := ctx.Value("status").(int); ok {
		r.AddAttrs(slog.Int("status", status))

	}
	return handler.Handler.Handle(ctx, r)
}
