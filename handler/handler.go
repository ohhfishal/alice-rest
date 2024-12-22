package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ohhfishal/alice-rest/config"
	"io"
	"log/slog"
	"net/http"
)

type Handler struct {
	Logger *slog.Logger
	Config *config.Config
}

type CustomHandler func(http.ResponseWriter, *http.Request) http.Handler

func (handler CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if next := handler(w, r); handler != nil {
		next.ServeHTTP(w, r)
		return
	}
	Text(400, "OK").ServeHTTP(w, r)
}

var ErrNotImplemented = errors.New("not implemented")

func status(err error) int {
	return 500
}

func JSON(status int, v any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(v)
	})
}

func Text(status int, content any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		fmt.Fprint(w, content)
	})
}

func Error(status int, err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, err.Error(), status)
	})
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	err := json.NewDecoder(r.Body).Decode(&v)
	if errors.Is(err, io.EOF) {
		return v, errors.New("body is empty or incomplete")
	}

	if err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}
