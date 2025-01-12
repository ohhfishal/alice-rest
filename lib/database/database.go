package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/ohhfishal/alice-rest/lib/event"
	"os"
	"sync"
)

var ErrNotFound = errors.New("not found")

type Database struct {
	mux       sync.Mutex
	directory string
}

type Option func(*Database) error

func New(options ...Option) (Database, error) {
	database := Database{}
	for _, option := range options {
		if err := option(&database); err != nil {
			return database, err
		}
	}
	return database, nil
}

func MountDirectory(path string) Option {
	return func(database *Database) error {
		database.directory = path
		return nil
	}
}

func (database Database) Close(_ context.Context) error {
	return nil
}

type Filter string

func (database Database) filePath(user string) string {
	return fmt.Sprintf("%s%s-events.json", database.directory, user)

}

func (database Database) Create(user string, newEvent event.Event) (string, error) {
	database.mux.Lock()
	defer database.mux.Unlock()

	path := database.filePath(user)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return "opening file:", err
	}
	defer file.Close()

	err = newEvent.To(file)
	if err != nil {
		return "writing to file:", err
	}
	return "ID_NOT_IMPLEMENTED", nil

}

func (database Database) Update(user, id string, update event.Event) (event.Event, error) {
	var temp event.Event
	database.mux.Lock()
	defer database.mux.Unlock()
	return temp, errors.New("not implemented")
}

func (database Database) Delete(user string, id string) error {
	database.mux.Lock()
	defer database.mux.Unlock()
	return errors.New("Not implemented")
	events, err := database.list(user)
	if err != nil {
		return fmt.Errorf("fetching events: %w", err)
	}

	path := database.filePath(user)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	found := false
	for _, event := range events {
		if event.ID == id {
			found = true
			continue
		}
		event.To(file)
	}

	if !found {
		return ErrNotFound
	}
	return nil

}

func (database Database) Get(user string, id string) (event.Event, error) {
	var temp event.Event

	database.mux.Lock()
	defer database.mux.Unlock()

	return temp, errors.New("not implemented")
}

func (database Database) List(user string, filters ...Filter) ([]event.Event, error) {
	database.mux.Lock()
	defer database.mux.Unlock()

	return database.list(user, filters...)
}

// Only call when you have the lock
func (database Database) list(user string, filters ...Filter) ([]event.Event, error) {
	file, err := os.Open(database.filePath(user))
	if err != nil {
		return []event.Event{}, err
	}
	return event.NewFrom(file)
}
