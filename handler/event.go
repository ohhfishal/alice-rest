package handler

import (
	"net/http"
)

func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Implemented", 501)
}

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Implemented", 501)
}

func (h *Handler) PutEvent(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Implemented", 501)
}

func (h *Handler) PatchEvent(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Implemented", 501)
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Implemented", 501)
}
