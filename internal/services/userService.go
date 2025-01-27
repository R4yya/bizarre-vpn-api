package services

import (
	"fmt"

	"bizarre-vpn-api/internal/storage/models"
)

type UserService struct {
	storage UserStorage
}

// GetUser gets the user by Telegram ID through the repository
func (s *UserService) GetUserById(ID int64) (*models.BaseUser, error) {
	user, err := s.storage.GetUserById(ID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// RegisterUser registers the user if it does not already exist
// func (s *UserService) RegisterUser(user *models.FullUser) (int64, error) {
// 	existingUser, err := s.storage.GetUserById()(user.TelegramID)
// 	if err == nil && existingUser != nil {
// 		return 0, intStorage.ErrUserAlreadyExists
// 	}

// 	userID, err := s.storage.CreateUser(user)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to register user: %w", err)
// 	}

// 	return userID, nil
// }
