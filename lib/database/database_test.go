package database

import (
	"io"
	// "testing"
)

func CreateTo(writer io.WriteCloser) Option {
	return func(database *Database) error {
		database.create = func(string) (io.WriteCloser, error) {
			return writer, nil
		}
		return nil
	}
}

func ReadFrom(reader io.ReadCloser) Option {
	return func(database *Database) error {
		database.openRead = func(string) (io.ReadCloser, error) {
			return reader, nil
		}
		return nil
	}
}

func WriteTo(reader io.WriteCloser) Option {
	return func(database *Database) error {
		database.openWrite = func(string) (io.WriteCloser, error) {
			return reader, nil
		}
		return nil
	}
}
