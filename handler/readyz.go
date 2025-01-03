package handler

import (
	"net/http"
)

func (h *Handler) Readyz() http.Handler {
	return CustomHandler(func(w http.ResponseWriter, r *http.Request) http.Handler {
		return Text(200, "Ready")
	})
}
