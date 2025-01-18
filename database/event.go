package database

import (
	"context"
  _ "embed"
	"fmt"
)

type Event struct {
	ID          int64
	Description string `json:"description"`
	State       string
}

//go:embed sql/migration.sql
var migration string

//go:embed sql/insert_event.sql
var insertEvent string

//go:embed sql/select_event.sql
var selectEvent string

//go:embed sql/delete_event.sql
var deleteEvent string

func InsertEvent(ctx context.Context, db Database, event Event) (int64, error) {
	result, err := db.ExecContext(ctx, insertEvent, event.Description)
	if err != nil {
		return -1, fmt.Errorf("insert: %w", err)
	}
	return result.LastInsertId()
}


func SelectEvent(ctx context.Context, db Database, id int64) (*Event, error) {
	row := db.QueryRowContext(ctx, selectEvent, id)

	event := &Event{}
	err := row.Scan(&event.ID, &event.Description, &event.State)
	if err != nil {
		return nil, err
	}

	return event, nil

}

func DeleteEvent(ctx context.Context, db Database, id int64) (int64, error) {
	result, err := db.ExecContext(ctx, deleteEvent, id)
	if err != nil {
		return -1, fmt.Errorf("delete: %w", err)
	}
	return result.RowsAffected()
}
