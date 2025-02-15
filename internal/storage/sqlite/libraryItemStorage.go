package sqlite

import (
	"bizarre-vpn-api/internal/storage"
	"bizarre-vpn-api/internal/storage/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type LibraryItemStorage struct {
	db        Database
	tableName string
	items     *[]models.LibraryItem
}

func NewLibraryItemStorage(db Database, tableName string, items *[]models.LibraryItem) *LibraryItemStorage {
	s := &LibraryItemStorage{
		db:        db,
		tableName: tableName,
		items:     items,
	}

	s.mustInit()

	return s
}

func (s *LibraryItemStorage) mustInit() {
	op := fmt.Sprintf("storage.LibraryItemStorage.initWithTransaction(%v)", s.tableName)

	tx, err := s.db.Beginx()

	if err != nil {
		panic(fmt.Errorf("%v: failed to create transaction, %w", op, err))
	}

	err = s.initWithTransaction(tx)

	if err != nil {
		rollbackErr := tx.Rollback()

		if rollbackErr != nil {
			panic(fmt.Errorf("%v: rollback error: %w", op, rollbackErr))
		}

		panic(fmt.Errorf("%v: %w", op, err))
	}

	tx.Commit()
}

func (s *LibraryItemStorage) initWithTransaction(executor storage.Executor) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %v(
		id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT NULL
	)  WITHOUT ROWID;`, s.tableName)

	_, err := executor.Exec(query)

	if err != nil {
		return fmt.Errorf("failed to init %v table: %w", s.tableName, err)
	}

	insertQuery := fmt.Sprintf(`INSERT OR IGNORE INTO %v (id, name, description) VALUES 
	(:id, :name, :description)`, s.tableName)

	_, err = sqlx.NamedExec(executor, insertQuery, *s.items)

	if err != nil {
		return fmt.Errorf("failed to fill rows of %v table: %w", s.tableName, err)
	}

	return nil
}

func (s *LibraryItemStorage) GetList() (*[]models.LibraryItem, error) {
	query := fmt.Sprintf(`SELECT * FROM %v`, s.tableName)

	var list []models.LibraryItem

	err := s.db.Select(&list, query)

	if err != nil {
		return nil, err
	}

	return &list, nil
}
