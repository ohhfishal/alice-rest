package server

import (
	"github.com/ohhfishal/alice-rest/handler"
	"net/http"
)

func addRoutes(mux *http.ServeMux, h *handler.Handler) http.Handler {

	mux.Handle("GET /api/v1/event", http.HandlerFunc(h.GetEvents))
	mux.Handle("GET /api/v1/event/{user}/{id}", http.HandlerFunc(h.GetEvent))
	mux.Handle("PUT /api/v1/event/{id}", http.HandlerFunc(h.PutEvent))
	mux.Handle("PATCH /api/v1/event/{id}", http.HandlerFunc(h.PatchEvent))
	mux.Handle("DELETE /api/v1/event/{id}", http.HandlerFunc(h.DeleteEvent))

	mux.HandleFunc("/", http.NotFound)

	return h.Log(mux)
}
