package event_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/ohhfishal/alice-rest/lib/event"
)

func cmp(e1, e2 event.Event) error {
	switch {
	case e1.Description != e2.Description:
		return fmt.Errorf("Description not equal %v != %v", e1.Description, e2.Description)
	case !e1.Due.Equal(e2.Due):
		return fmt.Errorf("Results not equal %v != %v", e2.Due, e2.Due)
	default:
		return nil
	}
}

func TestMarshaling(t *testing.T) {
	tests := []struct {
		Description string
		Options     []event.Option
	}{
		{Description: "a", Options: []event.Option{}},
		{Description: "b", Options: []event.Option{event.Due(time.Now())}},
	}

	for _, test := range tests {
		e, err := event.New(test.Description, test.Options...)
		if err != nil {
			t.Errorf("New(%s, %v) got an error %v", test.Description, test.Options, err)
		}

		if e.Description != test.Description {
			t.Errorf("Description incorrect. Got %s expected %s", e.Description, test.Description)
		}

		buffer := bytes.NewBuffer([]byte{})
		e.To(buffer)
		result, err := event.NewFrom(buffer)
		if err != nil {
			t.Errorf("NewFrom got an error, %v", err)
		}

		if len(result) != 1 {
			t.Error("Did not get 1 event")
		}

		if err := cmp(*e, result[0]); err != nil {
			t.Error(err)
		}
	}
}
