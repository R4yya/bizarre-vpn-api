package services

import (
	"fmt"

	intStorage "bizarre-vpn-api/internal/storage"
	"bizarre-vpn-api/internal/storage/models"
)

type UserStorage interface {
	GetUserByTelegramID(telegramID int64) (*models.User, error)
	CreateUser(user *models.User) (int64, error)
}

type UserService struct {
	storage UserStorage
}

// GetUser gets the user by Telegram ID through the repository
func (s *UserService) GetUser(telegramID int64) (*models.User, error) {
	user, err := s.storage.GetUserByTelegramID(telegramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// RegisterUser registers the user if it does not already exist
func (s *UserService) RegisterUser(user *models.User) (int64, error) {
	existingUser, err := s.storage.GetUserByTelegramID(user.TelegramID)
	if err == nil && existingUser != nil {
		return 0, intStorage.ErrUserAlreadyExists
	}

	userID, err := s.storage.CreateUser(user)
	if err != nil {
		return 0, fmt.Errorf("failed to register user: %w", err)
	}

	return userID, nil
}
