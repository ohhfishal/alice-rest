package alice

import (
	"context"
	alice "github.com/ohhfishal/alice-rest/lib/database"
	"github.com/ohhfishal/alice-rest/lib/event"
)

var MountDirectory = alice.MountDirectory

type Alice interface {
	// Returns ID, error
	Create(user string, newEvent event.Event) (string, error)
	Delete(user string, id string) error
	Get(user string, id string) (event.Event, error)
	List(string, ...alice.Filter) ([]event.Event, error)
	Close(context.Context) error
}

func New(args ...alice.Option) (Alice, error) {
	return alice.New(args...)
}
