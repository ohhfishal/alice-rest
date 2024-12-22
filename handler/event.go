package handler

import (
	"net/http"
	"time"
)

type Event struct {
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
}

func (h *Handler) GetEvent() http.Handler {
	return CustomHandler(func(w http.ResponseWriter, r *http.Request) http.Handler {
		_, err := decode[Event](r)
		if err != nil {
			return Error(400, err)
		}
		return Error(501, ErrNotImplemented)
	})
}

func (h *Handler) GetEvents() http.Handler {
	return CustomHandler(func(http.ResponseWriter, *http.Request) http.Handler {
		return Error(501, ErrNotImplemented)
	})
}

func (h *Handler) PostEvent() http.Handler {
	return CustomHandler(func(http.ResponseWriter, *http.Request) http.Handler {
		return Error(501, ErrNotImplemented)
	})
}

func (h *Handler) PatchEvent() http.Handler {
	return CustomHandler(func(http.ResponseWriter, *http.Request) http.Handler {
		return Error(501, ErrNotImplemented)
	})
}

func (h *Handler) DeleteEvent() http.Handler {
	return CustomHandler(func(http.ResponseWriter, *http.Request) http.Handler {
		return Error(501, ErrNotImplemented)
	})
}
