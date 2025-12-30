package services

import (
	"errors"
	"fmt"
	"log/slog"

	"bizarre-vpn-api/internal/lib/logger/sl"
	"bizarre-vpn-api/internal/storage/models"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorUserServiceInvalidPassword = errors.New("user password is invalid")
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
	const op = "internal.services.userService.GetUsersList"

	usersList, err := s.userStorage.GetUsersList()

	if err != nil {
		return nil, fmt.Errorf("%v: failed to get users list: %w", op, err)
	}

	return usersList, nil
}

func (s *UserService) GetUserById(ID int64) (*models.BaseUser, error) {
	const op = "internal.services.userService.GetUserById"

	user, err := s.userStorage.GetUserById(ID, nil)
	if err != nil {
		return nil, fmt.Errorf("%v: failed to get user: %w", op, err)
	}
	return user, nil
}

func (s *UserService) CreateUser(payload *models.CreateUserPayload) (*models.BaseUser, error) {
	const op = "internal.services.userService.CreateUser"

	if payload.Password != "" {
		//TODO: need to validate input password with error ErrorUserServiceInvalidPassword

		// Хэшируем пароль
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

		if err != nil {
			return nil, fmt.Errorf("%v: hashing password error: %w", op, err)
		}

		payload.Password = string(passwordHash)
	}

	user, err := s.userStorage.CreateUser(payload, nil)

	if err != nil {
		return nil, fmt.Errorf("%v: create user error: %w", op, err)
	}

	return user, nil
}

func (s *UserService) CreateDefaultUser() {
	const op = "internal.services.userService.CreateDefaultUser"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("Check of default user is exist")

	usersList, err := s.userStorage.GetUsersList()

	if err != nil {
		log.Error("read users list error", sl.Err(err))
		return
	}

	if len(*usersList) != 0 {
		return
	}

	defaultUserPayload := &models.CreateUserPayload{
		Username: "admin",
		Password: "admin",
		Role:     models.UserRoleAdmin,
	}

	_, err = s.CreateUser(defaultUserPayload)

	if err != nil {
		log.Error("create user error", sl.Err(err))
		return
	}

	log.Info(`Default user is created: username "admin"`)
}
