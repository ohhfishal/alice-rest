package server

import (
	"github.com/ohhfishal/alice-rest/handler"
	"net/http"
)

func addRoutes(mux *http.ServeMux, h *handler.Handler) http.Handler {

	mux.Handle("GET /api/v1/event", h.GetEvents())
	mux.Handle("POST /api/v1/event/{user}", h.PostEvent())
	mux.Handle("GET /api/v1/event/{user}/{id}", h.GetEvent())
	mux.Handle("PATCH /api/v1/event/{user}/{id}", h.PatchEvent())
	mux.Handle("DELETE /api/v1/event/{user}/{id}", h.DeleteEvent())

	mux.HandleFunc("/", http.NotFound)

	return h.Log(mux)
}
