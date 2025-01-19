package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type Storage struct {
	db     *sqlx.DB
	dbPath string
}

func New(dbPath string) (*Storage, error) {
	const op = "internal.storage.New"

	db, err := initDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", op, err)
	}

	storage := &Storage{
		db:     db,
		dbPath: dbPath,
	}

	return storage, nil
}

// InitDB initializes a connection to a database
func initDB(dbPath string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// CloseDB closes the connection to the database
func (s *Storage) CloseDB() error {
	if s.db == nil {
		return nil
	}

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}
