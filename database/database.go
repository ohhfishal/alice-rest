package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"log/slog"
)

var ErrNotFound = sql.ErrNoRows

type Database interface {
	Ping() error
	Exec(string, ...any) (sql.Result, error)
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRow(string, ...any) *sql.Row
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type SQLite struct {
	db     *sql.DB
	logger *slog.Logger
}

type Option func(*SQLite) error

func New(connectionString string, options ...Option) (Database, error) {
	db, err := sql.Open("sqlite", connectionString)
	if err != nil {
		return nil, fmt.Errorf("connection: %w", err)
	}

	sqlite := &SQLite{
		db: db,
		// TODO: Replace with slog.DiscardHandler in go 1.24
		logger: slog.New(slog.NewTextHandler(nil, &slog.HandlerOptions{Level: slog.LevelError})),
	}

	for _, option := range options {
		err := option(sqlite)
		if err != nil {
			return nil, fmt.Errorf("applying option: %w", err)
		}
	}

	// TODO: Move this out to the code that sets up the sql database
	_, err = sqlite.Exec(setup)
	if err != nil {
		return nil, fmt.Errorf("migrations: %w", err)
	}

	return sqlite, nil
}

func (s *SQLite) Ping() error {
	return s.db.Ping()
}

func (s *SQLite) Exec(query string, args ...any) (sql.Result, error) {
	return s.ExecContext(context.Background(), query, args...)
}

func (s *SQLite) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s *SQLite) QueryRow(query string, args ...any) *sql.Row {
	return s.QueryRowContext(context.Background(), query, args...)
}

func (s *SQLite) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}
