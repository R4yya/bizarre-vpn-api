package storage

import (
	"bizarre-vpn-api/pkg/logger"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var db *sqlx.DB

// InitDB initializes a connection to a database
func InitDB(dbPath string) error {
	var err error
	db, err = sqlx.Connect("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.Info(fmt.Sprintf("Connected to SQLite database at %s", dbPath))
	return nil
}

// GetDB returns the database object
func GetDB() *sqlx.DB {
	return db
}

// CloseDB closes the connection to the database
func CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			logger.Error(fmt.Errorf("failed to close database: %w", err))
		}
	}
}
