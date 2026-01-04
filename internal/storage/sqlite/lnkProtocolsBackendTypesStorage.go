package sqlite

import (
	"bizarre-vpn-api/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type LnkProtocolsBackendTypesStorage struct {
	db Database
}

func (ps *LnkProtocolsBackendTypesStorage) MustInit() {
	op := "storage.lnkProtocolsBackendTypesStorage.initWithTransaction"

	tx, err := ps.db.Beginx()

	if err != nil {
		panic(fmt.Errorf("%v: failed to create transaction, %w", op, err))
	}

	err = ps.initWithTransaction(tx)

	if err != nil {
		rollbackErr := tx.Rollback()

		if rollbackErr != nil {
			panic(fmt.Errorf("%v: rollback error: %w", op, rollbackErr))
		}

		panic(fmt.Errorf("%v: %w", op, err))
	}

	tx.Commit()
}

func (ps *LnkProtocolsBackendTypesStorage) initWithTransaction(executor Executor) error {
	query := `CREATE TABLE IF NOT EXISTS lnk_protocols_backend_types(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		protocol_id INTEGER NOT NULL,
		backend_type_id INTEGER NOT NULL,
		UNIQUE (protocol_id, backend_type_id),
    FOREIGN KEY (protocol_id) REFERENCES protocols(id),
    FOREIGN KEY (backend_type_id) REFERENCES backend_types(id)
	)`

	_, err := executor.Exec(query)

	if err != nil {
		return fmt.Errorf("failed to init lnk_protocols_backend_types table: %w", err)
	}

	insertQuery := `INSERT OR IGNORE INTO lnk_protocols_backend_types (protocol_id, backend_type_id) VALUES 
	(:protocol_id, :backend_type_id)`

	items := []models.LnkProtocolsBackendTypes{
		{
			ProtocolId:    1,
			BackendTypeId: 1,
		},
	}

	_, err = sqlx.NamedExec(executor, insertQuery, items)

	if err != nil {
		return fmt.Errorf("failed to fill rows of lnk_protocols_backend_types table: %w", err)
	}

	return nil
}

func (s *LnkProtocolsBackendTypesStorage) GetProtocolsByBackendTypeId(backendTypeId int64) (*[]models.LibraryItem, error) {
	query := `SELECT p.id AS id, p.name AS name
		FROM lnk_protocols_backend_types lpbt
		JOIN protocols p ON lpbt.protocol_id = p.id
		WHERE lpbt.backend_type_id = ?
		`

	items := make([]models.LibraryItem, 0)

	err := s.db.Select(&items, query, backendTypeId)

	if err != nil {
		return nil, fmt.Errorf("error with getting protocols by backend type id %w", err)
	}

	return &items, nil
}
