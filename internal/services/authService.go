package services

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"bizarre-vpn-api/internal/lib/jwt"
	"bizarre-vpn-api/internal/lib/logger/sl"
	"bizarre-vpn-api/internal/storage"
	"bizarre-vpn-api/internal/storage/models"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorAuthServiceIncorrectUsernameOrPass = errors.New("incorrect username or password")
)

const telegramProviderName string = "telegram"

type AuthService struct {
	log                    *slog.Logger
	userStorage            UserStorage
	lnkUserProviderStorage LnkUserProviderStorage
}

func NewAuthService(
	log *slog.Logger,
	userStorage UserStorage,
	lnkUserProviderStorage LnkUserProviderStorage,
) *AuthService {
	return &AuthService{
		userStorage:            userStorage,
		lnkUserProviderStorage: lnkUserProviderStorage,
		log:                    log,
	}
}

func (au *AuthService) getAuthorizeTokens(
	user *models.BaseUser,
	accessSecretKey []byte,
	refreshSecretKey []byte,
) (accessToken string, refreshToken string, Err error) {
	op := "auth.Authorize"

	log := au.log.With(slog.String("op", op))

	log.Debug("start generate tokens")

	accessToken, refreshToken, err := jwt.GenerateAuthTokens(
		user.ID,
		user.Role,
		accessSecretKey,
		refreshSecretKey,
	)

	if err != nil {
		log.Debug("generate tokens error", sl.Err(err))

		return "", "", fmt.Errorf("%v: %w", op, err)
	}

	log.Debug("generate tokens successful")

	log.Debug("update user refresh toke start")
	err = au.userStorage.UpdateUserRefreshToken(user.ID, refreshToken)

	if err != nil {
		log.Debug("update user refresh token error", sl.Err(err))
		return "", "", fmt.Errorf("%v: %w", op, err)
	}

	log.Debug("update user refresh token successful")

	return accessToken, refreshToken, nil
}

func (au *AuthService) authorizeByProvider(
	providerType string,
	externalUserId string,
	username string,
	accessSecretKey []byte,
	refreshSecretKey []byte,
) (accessToken string, refreshToken string, Err error) {
	op := "internal.services.userService.auth.authorizeByProvider"

	log := au.log.With(slog.String("op", op))

	lnkUserProvider, isFound, err := au.lnkUserProviderStorage.GetItemByType(
		telegramProviderName,
		externalUserId,
	)

	log.Debug("lnkUserProviderStorage.GetItemByType",
		slog.Any("lnkUserProvider", lnkUserProvider),
		slog.Bool("isFound", isFound),
		slog.String("externalUserId", externalUserId),
	)

	if err != nil {
		return "", "", fmt.Errorf("%v: %w", op, err)
	}

	var user *models.BaseUser

	if !isFound {
		log.Debug("in not found logic")
		newLnkUserProvider := &models.LnkUserProvider{
			ProviderType:   providerType,
			ExternalUserId: externalUserId,
		}

		_, user, err = au.lnkUserProviderStorage.CreateLnkUserProvider(newLnkUserProvider, username)

		log.Debug("CreateLnkUserProvider", slog.Any("user", user))

		if err != nil {
			return "", "", fmt.Errorf("%v: %w", op, err)
		}
	} else {
		log.Debug("in success found logic")
		user, err = au.userStorage.GetUserById(lnkUserProvider.UserId, nil)

		if err != nil {
			return "", "", fmt.Errorf("%v: %w", op, err)
		}
	}

	log.Debug("authorize logic")

	return au.getAuthorizeTokens(
		user,
		accessSecretKey,
		refreshSecretKey,
	)
}

func (au *AuthService) AuthorizeByTelegram(
	telegramId int64,
	username string,
	accessSecretKey []byte,
	refreshSecretKey []byte,
) (accessToken string, refreshToken string, Err error) {
	preparedExternalId := strconv.Itoa(int(telegramId))

	return au.authorizeByProvider(
		telegramProviderName,
		preparedExternalId,
		username,
		accessSecretKey,
		refreshSecretKey,
	)
}

func (au *AuthService) AuthorizeByCredentials(
	login string,
	password string,
	accessSecretKey []byte,
	refreshSecretKey []byte,
) (accessToken string, refreshToken string, Err error) {
	op := "internal.services.auth.AuthorizeByCredentials"

	user, err := au.userStorage.GetUserByLogin(login)

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return "", "", ErrorAuthServiceIncorrectUsernameOrPass
		}

		return "", "", fmt.Errorf("%v: error when getting user: %w", op, err)
	}

	userHashedPassword, err := au.userStorage.GetUserPasswordHash(user.ID)

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return "", "", ErrorAuthServiceIncorrectUsernameOrPass
		}

		return "", "", fmt.Errorf("%v: error when getting user password hash : %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(userHashedPassword), []byte(password))

	if err != nil {
		return "", "", ErrorAuthServiceIncorrectUsernameOrPass
	}

	return au.getAuthorizeTokens(user, accessSecretKey, refreshSecretKey)
}

func (au *AuthService) GetRefreshedTokens(
	user *models.BaseUser,
	refreshToken string,
	accessSecretKey []byte,
	refreshSecretKey []byte,
) (newAccessToken string, newRefreshToken string, Err error) {
	op := "internal.services.auth.GetRefreshedTokens"
	actualRefreshToken, err := au.userStorage.GetUserRefreshToken(user.ID)

	if err != nil {
		return "", "", fmt.Errorf("%v: error when getting user actual token: %w", op, err)
	}

	if actualRefreshToken != refreshToken {
		return "", "", fmt.Errorf("%v: user actual refresh token is not equal to transferred refresh token", op)
	}

	return au.getAuthorizeTokens(user, accessSecretKey, refreshSecretKey)
}
