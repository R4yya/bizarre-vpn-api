package services

import (
	"fmt"
	"log/slog"
	"strconv"

	"bizarre-vpn-api/internal/lib/jwt"
	"bizarre-vpn-api/internal/lib/logger/sl"
	"bizarre-vpn-api/internal/storage/models"
)

const telegramProviderName string = "telegram"

type LnkUserProviderStorage interface {
	GetItemByType(
		providerType string,
		externalId string,
	) (LnkUserProvider *models.LnkUserProvider, isFound bool, Err error)
	CreateLnkUserProvider(
		lnkUserProvider *models.LnkUserProvider,
		username string,
	) (providerId int64, createdUser *models.BaseUser, Err error)
}

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
) (accessToken string, refreshToken string, Err error) {
	op := "auth.Authorize"

	log := au.log.With(slog.String("op", op))

	log.Debug("start generate tokens")

	accessToken, refreshToken, err := jwt.GenerateAuthTokens(user.ID, user.Role)

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

	return au.getAuthorizeTokens(user)
}

func (au *AuthService) AuthorizeByTelegram(
	telegramId int64,
	username string,
) (accessToken string, refreshToken string, Err error) {
	preparedExternalId := strconv.Itoa(int(telegramId))

	return au.authorizeByProvider(telegramProviderName, preparedExternalId, username)
}
