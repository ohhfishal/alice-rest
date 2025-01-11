package server

import (
	"github.com/ohhfishal/alice-rest/server/handler"
	"net/http"
)

func addRoutes(mux *http.ServeMux, h handler.Handler) http.Handler {

	mux.Handle("GET /api/v1/event/{user}", h.UserAuth(h.GetEvents()))
	mux.Handle("POST /api/v1/event/{user}", h.UserAuth(h.PostEvent()))
	mux.Handle("GET /api/v1/event/{user}/{id}", h.UserAuth(h.GetEvent()))
	mux.Handle("PATCH /api/v1/event/{user}/{id}", h.UserAuth(h.PatchEvent()))
	mux.Handle("DELETE /api/v1/event/{user}/{id}", h.UserAuth(h.DeleteEvent()))

	mux.Handle("GET /readyz", h.Readyz())

	mux.HandleFunc("/", http.NotFound)

	fullHandler := h.Log(mux)
	fullHandler = h.SetupContext(fullHandler)
	fullHandler = http.TimeoutHandler(fullHandler, h.ResponseTimeout, "Request Timeout")
	return fullHandler
}
