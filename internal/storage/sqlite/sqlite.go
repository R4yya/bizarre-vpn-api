package sqlite

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type Database = *sqlx.DB

type Storage struct {
	db               Database
	subscriptionPlan SubscriptionPlan
}

func Init(dbPath string, log *slog.Logger) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("%v: failed to connect to database: %w", op, err)
	}

	log.Info(fmt.Sprintf("Connected to SQLite database at %s", dbPath))

	storage := &Storage{
		db:               db,
		subscriptionPlan: SubscriptionPlan{db: db},
	}

	return storage, nil
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
