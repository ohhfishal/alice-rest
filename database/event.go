package database

import (
	"context"
	"fmt"
)

type Event struct {
	ID          int64
	Description string `json:"description"`
	State       string
}

// TODO: Move into a migration script
const setup = `
CREATE TABLE events (
  id INTEGER PRIMARY KEY AUTOINCREMENT, --sqllite bug
  description TEXT NOT NULL,
  state VARCHAR(20) NOT NULL DEFAULT 'in progress' CHECK (state IN ('in progress', 'done'))
);
`

const insertEvent = `INSERT INTO events (description) VALUES (?);`

func InsertEvent(ctx context.Context, db Database, event Event) (int64, error) {
	result, err := db.ExecContext(ctx, insertEvent, event.Description)
	if err != nil {
		return -1, fmt.Errorf("insert: %w", err)
	}
	return result.LastInsertId()
}

const getEvent = `SELECT * FROM events WHERE id = ?;`

func GetEvent(ctx context.Context, db Database, id int64) (*Event, error) {
	row := db.QueryRowContext(ctx, getEvent, id)

	event := &Event{}
	err := row.Scan(&event.ID, &event.Description, &event.State)
	if err != nil {
		return nil, err
	}

	return event, nil

}
