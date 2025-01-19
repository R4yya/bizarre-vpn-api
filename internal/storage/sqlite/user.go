package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"bizarre-vpn-api/internal/storage"
	"bizarre-vpn-api/internal/storage/models"
)

type User struct {
	db Database
}

// CreateUser adds a new user to the database
func (u *User) CreateUser(user *models.User) (int64, error) {
	query := `
    INSERT INTO users (telegram_id, username, language_code, is_bot)
    VALUES (:telegram_id, :username, :language_code, :is_bot)
    `
	result, err := u.db.NamedExec(query, user)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	return userID, nil
}

// GetUserByTelegramID gets the user by Telegram ID
func (u *User) GetUserByTelegramID(telegramID int64) (*models.User, error) {
	var user models.User

	query := "SELECT * FROM users WHERE telegram_id = ?"
	err := u.db.Get(&user, query, telegramID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}
