package server

import (
	"github.com/ohhfishal/alice-rest/server/handler"
	"log/slog"
	"net/http"
)

func NewServer(h *handler.Handler) http.Handler {
	if h == nil {
		slog.Error("[NewServer] Handler is nil")
	}
	mux := http.NewServeMux()
	return addRoutes(mux, h)
}
