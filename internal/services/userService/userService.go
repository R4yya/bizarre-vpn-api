package userService

import (
	"errors"
	"fmt"
	"log/slog"
	"unicode/utf8"

	"bizarre-vpn-api/internal/lib/logger/sl"
	"bizarre-vpn-api/internal/services"
	"bizarre-vpn-api/internal/storage"
	"bizarre-vpn-api/internal/storage/models"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("user password is invalid")
	ErrIncorrectRole   = errors.New("role is incorrect")
	ErrLoginIsTooSmall = errors.New("login is too small")
	ErrLoginOccupied   = errors.New("user with this login already exist")
)

type UserService struct {
	log         *slog.Logger
	userStorage services.UserStorage
}

func NewUserService(
	log *slog.Logger,
	userStorage services.UserStorage,
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

	log := s.log.With(
		slog.String("op", op),
	)

	log.Debug("CreateUser start")

	if payload.Role != models.UserRoleClient && payload.Role != models.UserRoleAdmin {
		return nil, ErrIncorrectRole
	}

	if payload.Login != nil {
		if utf8.RuneCountInString(*payload.Login) < 5 {
			return nil, ErrLoginIsTooSmall
		}
	}

	if payload.Password != nil {
		//TODO: add more strange password validation
		if utf8.RuneCountInString(*payload.Password) < 5 {
			return nil, ErrInvalidPassword
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(*payload.Password), bcrypt.DefaultCost)

		if err != nil {
			return nil, fmt.Errorf("%v: hashing password error: %w", op, err)
		}

		stringHash := string(passwordHash)

		payload.Password = &stringHash
	}

	user, err := s.userStorage.CreateUser(payload, nil)

	if err != nil {
		if errors.Is(err, storage.ErrLoginOccupied) {
			return nil, ErrLoginOccupied
		}

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

	log.Debug("len(*usersList)", slog.Any("len(*usersList)", len(*usersList)))

	if len(*usersList) != 0 {
		return
	}

	defaultUserLogin := "admin"

	defaultUserPayload := &models.CreateUserPayload{
		Username: "Default admin user",
		Login:    &defaultUserLogin,
		Password: &defaultUserLogin,
		Role:     models.UserRoleAdmin,
	}

	_, err = s.CreateUser(defaultUserPayload)

	if err != nil {
		log.Error("create user error", sl.Err(err))
		return
	}

	log.Info(`Default user is created: username "admin"`)
}
