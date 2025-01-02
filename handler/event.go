package handler

import (
	"errors"
	"fmt"
	"net/http"

	alice "github.com/ohhfishal/alice-rest/lib/event"
)

func (h *Handler) GetEvent() http.Handler {
	return CustomHandler(func(w http.ResponseWriter, r *http.Request) http.Handler {
		user := r.PathValue("user")
		if user == "" {
			return Error(400, errors.New("missing user"))
		}

		id := r.PathValue("id")
		if user == "" {
			return Error(400, errors.New("id"))
		}

		event, err := h.Alice.Get(user, id)
		if err != nil {
			return Error(500, fmt.Errorf("getting event: %w", err))
		}
		return JSON(200, event)
	})
}

func (h *Handler) GetEvents() http.Handler {
	return CustomHandler(func(w http.ResponseWriter, r *http.Request) http.Handler {
		return Error(501, ErrNotImplemented)
	})
}

func (h *Handler) PostEvent() http.Handler {
	return CustomHandler(func(w http.ResponseWriter, r *http.Request) http.Handler {
		user := r.PathValue("user")
		if user == "" {
			return Error(400, errors.New("missing user"))
		}

		newEvent, err := decode[alice.Event](r)
		if err != nil {
			return Error(400, err)
		}

		_, err = h.Alice.Create(user, newEvent)
		if err != nil {
			return Error(500, fmt.Errorf("failed to create event: %w", err))
		}
		return Text(201, "Created\n")
	})
}

func (h *Handler) PatchEvent() http.Handler {
	return CustomHandler(func(w http.ResponseWriter, r *http.Request) http.Handler {
		return Error(501, ErrNotImplemented)
	})
}

func (h *Handler) DeleteEvent() http.Handler {
	return CustomHandler(func(w http.ResponseWriter, r *http.Request) http.Handler {
		user := r.PathValue("user")
		if user == "" {
			return Error(400, errors.New("missing user"))
		}

		id := r.PathValue("id")
		if user == "" {
			return Error(400, errors.New("id"))
		}

		err := h.Alice.Delete(user, id)
		if err != nil {
			return Error(500, fmt.Errorf("deleting event: %w", err))
		}
		return Text(200, "Deleted\n")
	})
}
