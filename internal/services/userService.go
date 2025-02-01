package services

import (
	"fmt"
	"log/slog"

	"bizarre-vpn-api/internal/storage/models"
)

type UserService struct {
	userStorage UserStorage
}

func NewUserService(
	log *slog.Logger,
	userStorage UserStorage,
) *UserService {
	return &UserService{
		userStorage: userStorage,
	}
}

// GetUser gets the user by Telegram ID through the repository
func (s *UserService) GetUserById(ID int64) (*models.BaseUser, error) {
	user, err := s.userStorage.GetUserById(ID, nil)
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
