package services

import (
	"fmt"
	"log/slog"

	"bizarre-vpn-api/internal/storage/models"
)

const (
	UserBasicRole string = "basic"
	UserAdminRole string = "admin"
)

type UserService struct {
	log         *slog.Logger
	userStorage UserStorage
}

func NewUserService(
	log *slog.Logger,
	userStorage UserStorage,
) *UserService {
	return &UserService{
		log:         log,
		userStorage: userStorage,
	}
}

func (s *UserService) GetUsersList() (*[]models.BaseUser, error) {
	const op = "internal.services.GetUsersList"

	usersList, err := s.userStorage.GetUsersList()

	if err != nil {
		return nil, fmt.Errorf("%v: failed to get users list: %w", op, err)
	}

	return usersList, nil
}

func (s *UserService) GetUserById(ID int64) (*models.BaseUser, error) {
	const op = "internal.services.GetUserById"

	user, err := s.userStorage.GetUserById(ID, nil)
	if err != nil {
		return nil, fmt.Errorf("%v: failed to get user: %w", op, err)
	}
	return user, nil
}
