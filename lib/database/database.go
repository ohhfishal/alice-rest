package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/ohhfishal/alice-rest/lib/event"
	"io"
	"os"
	"sync"
)

var ErrNotFound = errors.New("not found")

type Database struct {
	mux       sync.Mutex
	create    func(string) (io.WriteCloser, error)
	openRead  func(string) (io.ReadCloser, error)
	openWrite func(string) (io.WriteCloser, error)
}

type Option func(*Database) error

func New(options ...Option) (Database, error) {
	database := Database{
		create:    create,
		openRead:  openRead,
		openWrite: openWrite,
	}
	for _, option := range options {
		if err := option(&database); err != nil {
			return database, err
		}
	}
	return database, nil
}

func (database Database) Close(_ context.Context) error {
	return nil
}

type Filter string

func (database Database) filePath(user string) string {
	return fmt.Sprintf("%s-events.json", user)

}

func (database Database) Create(user string, newEvent event.Event) (string, error) {
	database.mux.Lock()
	defer database.mux.Unlock()

	filePath := database.filePath(user)
	file, err := database.openWrite(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = newEvent.To(file)
	if err != nil {
		return "", err
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
	events, err := database.list(user)
	if err != nil {
		return fmt.Errorf("fetching events: %w", err)
	}

	file, err := database.create(database.filePath(user))
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
	file, err := database.openRead(database.filePath(user))
	if err != nil {
		return []event.Event{}, err
	}
	return event.NewFrom(file)
}

func create(name string) (io.WriteCloser, error) {
	return os.Create(name)
}

func openRead(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

func openWrite(name string) (io.WriteCloser, error) {
	return os.OpenFile(name, os.O_WRONLY|os.O_APPEND, 0777)
}
