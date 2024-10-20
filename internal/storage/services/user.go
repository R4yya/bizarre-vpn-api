package services

import (
	"bizarre-vpn-api/internal/storage/models"
	"bizarre-vpn-api/internal/storage/repositories"
	"fmt"
)

// GetUser gets the user by Telegram ID through the repository
func GetUser(telegramID int64) (*models.User, error) {
	user, err := repositories.GetUserByTelegramID(telegramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// RegisterUser registers the user if it does not already exist
func RegisterUser(user *models.User) (int64, error) {
	existingUser, err := repositories.GetUserByTelegramID(user.TelegramID)
	if err == nil && existingUser != nil {
		return 0, fmt.Errorf("user with Telegram ID %d already exists", user.TelegramID)
	}

	userID, err := repositories.CreateUser(user)
	if err != nil {
		return 0, fmt.Errorf("failed to register user: %w", err)
	}

	return userID, nil
}
