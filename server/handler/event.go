package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ohhfishal/alice-rest/database"
)

func (h *Handler) PostEvent() http.Handler {
	return CustomHandler(func(w http.ResponseWriter, r *http.Request) http.Handler {
		_ = r.PathValue("user")

		event, err := decode[database.Event](r)
		if err != nil {
			return Error(400, err)
		}

		// TODO: Add the context here
		id, err := database.InsertEvent(context.Background(), h.DB, event)
		if err != nil {
			return Error(500, fmt.Errorf("getting event: %w", err))
		}
		return JSON(201, id)
	})
}

func (h *Handler) GetEvent() http.Handler {
	return CustomHandler(func(w http.ResponseWriter, r *http.Request) http.Handler {
		_ = r.PathValue("user")
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return Error(400, err)
		}

		event, err := database.GetEvent(context.Background(), h.DB, (int64)(id))
		switch {
		case errors.Is(database.ErrNotFound, err):
			return Error(404, fmt.Errorf("create event: %w", err))
		case err != nil:
			return Error(500, fmt.Errorf("create event: %w", err))
		default:
			return JSON(200, event)
		}
	})
}

func (h *Handler) GetEvents() http.Handler {
	return CustomHandler(func(w http.ResponseWriter, r *http.Request) http.Handler {
		return Error(501, ErrNotImplemented)
	})
}

func (h *Handler) PatchEvent() http.Handler {
	return CustomHandler(func(w http.ResponseWriter, r *http.Request) http.Handler {
		return Error(501, ErrNotImplemented)
	})
}

func (h *Handler) DeleteEvent() http.Handler {
	return CustomHandler(func(w http.ResponseWriter, r *http.Request) http.Handler {
		return Error(501, ErrNotImplemented)
	})
}
