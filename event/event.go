package event

import (
	"errors"
	"fmt"
	alice "github.com/ohhfishal/alice/database"
	"github.com/ohhfishal/alice/event"
)

type EventManager interface {
	Create(user string, newEvent event.Event) (string, error)
	Delete(user string, id string) error
	Get(user string, id string) (event.Event, error)
	List(string, ...Filter) ([]event.Event, error)
}

func New() EventManager {
	return alice.New()
}
