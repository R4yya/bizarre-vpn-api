package sqlite

import (
	"bizarre-vpn-api/internal/core/coreErrors"
	"bizarre-vpn-api/internal/storage/models"
	"fmt"
)

type VpnServersStorage struct {
	db Database
}

func (s *VpnServersStorage) MustInit() {
	query := `CREATE TABLE IF NOT EXISTS vpn_servers(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		adapter_host TEXT NOT NULL,
		adapter_port INTEGER NOT NULL,
		country TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		lnk_backend_type_id INTEGER NOT NULL,
		FOREIGN KEY (lnk_backend_type_id) REFERENCES lnk_protocols_backend_types(id),
		UNIQUE(adapter_host, adapter_port)
	)`

	_, err := s.db.Exec(query)

	if err != nil {
		panic(fmt.Errorf("failed to init vpn_servers table: %w", err))
	}
}

func (s *VpnServersStorage) GetExpandedList() (*[]models.VpnServerExpandedItem, error) {
	query := `SELECT vse.*, p.name AS protocol, bt.name AS backend_type FROM vpn_servers vse
	JOIN lnk_protocols_backend_types lpbt ON vse.lnk_backend_type_id = lpbt.id
	JOIN protocols p ON lpbt.protocol_id = p.id
	JOIN backend_types bt ON lpbt.backend_type_id = bt.id
	`

	vpnServersList := []models.VpnServerExpandedItem{}

	err := s.db.Select(&vpnServersList, query)

	if err != nil {
		return nil, fmt.Errorf("failed to get vpn_servers list: %w", err)
	}

	return &vpnServersList, nil
}

func (s *VpnServersStorage) GetExpandedItemById(vpnServerId int64) (*models.VpnServerExpandedItem, error) {
	query := `SELECT vse.*, p.name AS protocol, bt.name AS backend_type FROM vpn_servers vse
	JOIN lnk_protocols_backend_types lpbt 
		ON vse.lnk_backend_type_id = lpbt.id
	JOIN protocols p 
		ON lpbt.protocol_id = p.id
	JOIN backend_types bt
		ON lpbt.backend_type_id = bt.id
	WHERE vse.id = ?
	`

	vpnServer := models.VpnServerExpandedItem{}

	err := s.db.Get(&vpnServer, query, vpnServerId)

	if err != nil {
		return nil, coreErrors.ErrorNotFound
	}

	fmt.Printf("\n vpnServer %v \n", vpnServer)

	return &vpnServer, nil
}

func (s *VpnServersStorage) CreateItem(vpnServer *models.VpnServerItem) (int64, error) {
	query := `INSERT INTO vpn_servers 
	(
		name,
		description,
		adapter_host,
		adapter_port,
		country,
		lnk_backend_type_id
	) 
	VALUES (
		:name,
		:description,
		:adapter_host,
		:adapter_port,
		:country,
		:lnk_backend_type_id
	)`

	result, err := s.db.NamedExec(query, vpnServer)

	if err != nil {
		return 0, fmt.Errorf("failed to create vpn_server: %v", err)
	}

	vpnServerID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	return vpnServerID, nil
}

func (s *VpnServersStorage) DeleteItem(vpnServerId int64) error {
	query := `DELETE FROM vpn_servers WHERE id = ?`

	result, err := s.db.Exec(query, vpnServerId)
	if err != nil {
		return fmt.Errorf("failed to delete vpn_server: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return coreErrors.ErrorNotFound
	}

	return nil
}
