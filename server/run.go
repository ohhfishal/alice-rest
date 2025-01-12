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

	"github.com/ohhfishal/alice-rest/lib/alice"
	"github.com/ohhfishal/alice-rest/server/handler"
)

func Run(
	ctx context.Context, args []string, getenv func(string) string, stdin io.Reader, stdout, stderr io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cfg := NewConfig(args, getenv)

	logger := NewLogger(stdout, cfg.LogLevel)

	a, err := alice.New(alice.MountDirectory(cfg.DatabaseDirectory))
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	h := handler.Handler{
		Logger:          logger,
		Alice:           a,
		ResponseTimeout: cfg.ResponseTimeout,
	}

	mux := http.NewServeMux()
	server := addRoutes(mux, h)
	httpServer := &http.Server{
		Addr:         net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:      server,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: cfg.ResponseTimeout + time.Second,
	}

	go func() {
		logger.Info(fmt.Sprintf("listening on %s", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error(fmt.Errorf("listening and serving: %w", err).Error())
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		logger.Info("shutting down")
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error(fmt.Errorf("shutting down: %w", err).Error())
		}
	}()
	wg.Wait()
	return nil
}
