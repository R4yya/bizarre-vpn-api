package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"bizarre-vpn-api/internal/storage"
	"bizarre-vpn-api/internal/storage/models"

	"github.com/jmoiron/sqlx"
)

type LnkUserProviderStorage struct {
	db Database
}

func (s *LnkUserProviderStorage) MustInit() {
	query := `CREATE TABLE IF NOT EXISTS lnk_user_providers(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		provider_type TEXT NOT NULL,
		external_user_id TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE (user_id, provider_type)
	)`

	_, err := s.db.Exec(query)

	if err != nil {
		panic(fmt.Errorf("failed to init lnk_user_providers table: %w", err))
	}
}

func (s *LnkUserProviderStorage) GetItemByType(
	providerType string,
	externalId string,
) (LnkUserProvider *models.LnkUserProvider, isFiend bool, Err error) {
	query := `SELECT * FROM lnk_user_providers WHERE provider_type = ? and external_user_id = ?`

	var lnkUserProvider = &models.LnkUserProvider{}

	err := s.db.Get(lnkUserProvider, query, providerType, externalId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}

		return nil, false, fmt.Errorf(
			"failed to getting LnkUserProvider FOR providerType=%v, external_user_id=%v, err:%w",
			providerType,
			externalId,
			err,
		)
	}

	return lnkUserProvider, true, nil
}

func (s *LnkUserProviderStorage) GetListByUserId(userId int64) (*[]models.LnkUserProvider, error) {
	query := `SELECT * FROM lnk_user_providers WHERE user_id = ?`

	var list []models.LnkUserProvider

	err := s.db.Select(&list, query, userId)

	if err != nil {
		return nil, fmt.Errorf("getting user providers error: %w", err)
	}

	return &list, nil
}

func (s *LnkUserProviderStorage) CreateLnkUserProvider(
	createLnkUserProviderPayload *models.CreateLnkUserProviderPayload,
	executor storage.Executor,
) (*models.LnkUserProvider, error) {
	if executor == nil {
		executor = s.db
	}

	query := `INSERT into lnk_user_providers 
	(provider_type, external_user_id, user_id) 
	VALUES (:provider_type, :external_user_id, :user_id) RETURNING *`

	rows, err := sqlx.NamedQuery(executor, query, createLnkUserProviderPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to create LnkUserProvider: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("failed to get rows.Next: %w", err)
	}

	var createdProvider models.LnkUserProvider

	err = rows.StructScan(&createdProvider)

	if err != nil {
		return nil, fmt.Errorf("failed to scan inserted LnkUserProvider: %w", err)
	}

	return &createdProvider, nil
}

// func (s *LnkUserProviderStorage) addProviderIfNotExist(authProvider *models.AuthProvider) (int64, error) {

// 	query := `
// 	INSERT INTO auth_providers (:name, :external_user_id_format),
// 	VALUES (:name, :external_user_id_format)
// 	`

// 	result, err := s.db.NamedExec(query, authProvider)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to create authProvider: %w", err)
// 	}

// 	userID, err := result.LastInsertId()
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to retrieve last insert ID: %w", err)
// 	}

// 	return userID, nil
// }
