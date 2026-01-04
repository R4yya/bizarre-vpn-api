package services

import (
	"errors"
	"fmt"
	"log/slog"
	"unicode/utf8"

	"bizarre-vpn-api/internal/models"
	"bizarre-vpn-api/internal/services/interfaces"
	"bizarre-vpn-api/internal/shared/coreErrors"
	"bizarre-vpn-api/internal/shared/logger/sl"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	log         *slog.Logger
	userStorage interfaces.UserStorage
}

func NewUserService(
	log *slog.Logger,
	userStorage interfaces.UserStorage,
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
		if errors.Is(err, coreErrors.ErrorNotFound) {
			return nil, err
		}

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
		return nil, coreErrors.ValidationError{
			Entity: "user",
			Msg:    "incorrect value",
			Fields: []string{
				"role",
			},
		}
	}

	if payload.Login != nil {
		if loginLength := utf8.RuneCountInString(*payload.Login); loginLength < 5 || loginLength >= 20 {
			return nil, coreErrors.ValidationError{
				Entity: "user",
				Msg:    "must be longer than 5 and shorter than 20",
				Fields: []string{
					"login",
				},
			}
		}
	}

	if payload.Password != nil {
		//TODO: add more strange password validation
		if passwordLength := utf8.RuneCountInString(*payload.Password); passwordLength < 5 || passwordLength >= 250 {
			return nil, coreErrors.ValidationError{
				Entity: "user",
				Msg:    "must be longer than 5 and shorter than 250",
				Fields: []string{
					"password",
				},
			}
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
		if errors.Is(err, coreErrors.ErrorAlreadyExist) {
			return nil, coreErrors.ErrorAlreadyExist
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
