package database

import (
	"context"
  "log/slog"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testContext() context.Context {
	return context.Background()
}

func newMemoryDB(t *testing.T) Database {
  logger := slog.Default()
	db, err := New(":memory:", Log(logger), Migrate(testContext()))
	require.Nil(t, err)
	require.NotNil(t, db)
	return db
}

func TestInit(t *testing.T) {
	db := newMemoryDB(t)
	assert.Nil(t, db.Ping(), "failed to ping")
}

func TestEvent(t *testing.T) {
	db := newMemoryDB(t)
	assert.Nil(t, db.Ping())

	event := Event{
		Description: "foo",
	}

	ctx := testContext()

	id, err := InsertEvent(ctx, db, event)
	require.Nil(t, err)
	t.Logf("ID: %v", id)

	result, err := db.ExecContext(ctx, "SELECT COUNT(*) from events WHERE id = 1")
	require.Nil(t, err)
	count, err := result.RowsAffected()
	require.Nil(t, err)
	t.Logf("Count: %d", count)

	newEvent, err := SelectEvent(ctx, db, id)
	require.Nil(t, err)

	// TODO: Make a function to assert events are the same
	assert.Equal(t, event.Description, newEvent.Description)
	assert.Equal(t, "in progress", newEvent.State)

  rows, err := DeleteEvent(ctx, db, id)
  require.Nil(t, err)
  require.Equal(t, (int64)(1), rows)

	_, err = SelectEvent(ctx, db, id)
	require.NotNil(t, err)
}
