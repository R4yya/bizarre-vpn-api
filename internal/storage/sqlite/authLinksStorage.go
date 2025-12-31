package sqlite

import (
	"bizarre-vpn-api/internal/storage"
	"bizarre-vpn-api/internal/storage/models"
	"fmt"
)

type AuthLinksStorage struct {
	db Database
}

func (s *AuthLinksStorage) MustInit() {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS auth_links(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		code TEXT NOT NULL UNIQUE,
		status TEXT NOT NULL DEFAULT %v,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	)`, models.AuthLinkStatusCreated)

	_, err := s.db.Exec(query)

	if err != nil {
		panic(fmt.Errorf("failed to init auth_links table: %w", err))
	}
}

func (s *AuthLinksStorage) GetItemByCode(code string) (*models.AuthLink, error) {
	query := `SELECT 
		id,
		user_id,
		code,
		status,
		created_at,
		updated_at 
		FROM auth_links WHERE code = ?`

	var authLink models.AuthLink

	err := s.db.Get(&authLink, query, code)

	if err != nil {
		return nil, storage.ErrAuthLinkNotFound
	}

	return &authLink, nil
}

func (s *AuthLinksStorage) CreateItem(payload *models.AuthLinkCreatePayload) (*models.AuthLink, error) {
	query := `INSERT INTO auth_links (user_id, code, status) VALUES (:user_id, :code, :status) RETURNING *`

	rows, err := s.db.NamedQuery(
		query,
		payload,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create auth link: %w", err)
	}

	fmt.Println("rows", rows)

	defer rows.Close()

	var authLink models.AuthLink

	if rows.Next() {
		err = rows.StructScan(&authLink)

		if err != nil {
			return nil, fmt.Errorf("failed to scan inserted auth link: %w", err)
		}
	} else {
		return nil, fmt.Errorf("not rows next")
	}

	return &authLink, nil
}
