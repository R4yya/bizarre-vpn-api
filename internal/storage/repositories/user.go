package repositories

import (
	"bizarre-vpn-api/internal/storage"
	"bizarre-vpn-api/internal/storage/models"
	"fmt"
)

// CreateUser adds a new user to the database
func CreateUser(user *models.User) (int64, error) {
	query := `
    INSERT INTO users (telegram_id, username, language_code, is_bot)
    VALUES (:telegram_id, :username, :language_code, :is_bot)
    `
	result, err := storage.GetDB().NamedExec(query, user)
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
func GetUserByTelegramID(telegramID int64) (*models.User, error) {
	var user models.User

	query := "SELECT * FROM users WHERE telegram_id = ?"
	err := storage.GetDB().Get(&user, query, telegramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}
