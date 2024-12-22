package server

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ohhfishal/alice-rest/config"
	"github.com/ohhfishal/alice-rest/event"
	"github.com/ohhfishal/alice-rest/handler"
)

func Run(
	ctx context.Context, args []string, getenv func(string) string, stdin io.Reader, stdout, stderr io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cfg := config.NewConfig(args, getenv)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	h := handler.Handler{
		Logger:       logger,
		Config:       cfg,
		EventManager: event.New(),
	}
	server := NewServer(&h)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: server,
	}

	go func() {
		logger.Info(fmt.Sprintf("listening on %s", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listening and serving: %s", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("shutting down: %w", err)
		}
	}()
	wg.Wait()
	return nil
}
