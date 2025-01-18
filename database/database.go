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
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type SQLite struct {
	db     *sql.DB
	logger *slog.Logger
}

type Option func(*SQLite) error

func Log(logger *slog.Logger) func(*SQLite) error {
  return func(s *SQLite) error {
    s.logger = logger
    return nil
  }
}

func Migrate(ctx context.Context) func(*SQLite) error {
  return func(s *SQLite) error {
    _, err := s.ExecContext(ctx, migration)
    if err != nil {
      return fmt.Errorf("migrations: %w", err)
    }
    return nil
  }
}

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
	return sqlite, nil
}

func (s *SQLite) Ping() error {
	return s.db.Ping()
}

func (s *SQLite) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
  s.logger.Debug("Exec", "sql", query, "args", args)
	return s.db.ExecContext(ctx, query, args...)
}

func (s *SQLite) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
  s.logger.Debug("Query", "sql", query, "args", args)
	return s.db.QueryRowContext(ctx, query, args...)
}
