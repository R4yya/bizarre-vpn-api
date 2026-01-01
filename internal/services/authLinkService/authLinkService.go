package authLinkService

import (
	"bizarre-vpn-api/internal/core/coreErrors"
	"bizarre-vpn-api/internal/lib/random"
	"bizarre-vpn-api/internal/services"
	"bizarre-vpn-api/internal/storage/models"
	"errors"
	"fmt"
	"log/slog"
)

const (
	LinkUserCodeLength = 6
)

type AuthLinksStorage interface {
	GetItemByCode(code string) (*models.AuthLink, error)
	CreateItem(payload *models.AuthLinkCreatePayload) (*models.AuthLink, error)
}

type AuthLinkService struct {
	log                    *slog.Logger
	authLinksStorage       AuthLinksStorage
	lnkUserProviderStorage services.LnkUserProviderStorage
}

func NewAuthLinksService(
	log *slog.Logger,
	authLinksStorage AuthLinksStorage,
	lnkUserProviderStorage services.LnkUserProviderStorage,
) *AuthLinkService {
	return &AuthLinkService{
		authLinksStorage:       authLinksStorage,
		log:                    log,
		lnkUserProviderStorage: lnkUserProviderStorage,
	}
}

func (s *AuthLinkService) LinkUserWithTgProviderByCode(code string, externalUserId string) (userId int64, Err error) {
	op := "internal.services.authLinksService.LinkUserWithTgProviderByCode"

	authLink, err := s.authLinksStorage.GetItemByCode(code)

	if err != nil {
		if errors.Is(err, coreErrors.ErrorNotFound) {
			return 0, err
		}

		return 0, fmt.Errorf("%v: %w", op, err)
	}

	userProviders, err := s.lnkUserProviderStorage.GetListByUserId(authLink.UserId)

	if err != nil {
		return 0, fmt.Errorf("%v: %w", op, err)
	}

	if len(*userProviders) != 0 {
		return 0, coreErrors.ErrorUserAlreadyLinked
	}

	createLnkUserProviderPayload := models.CreateLnkUserProviderPayload{
		ProviderType:   models.TelegramProviderName,
		ExternalUserId: externalUserId,
		UserId:         authLink.UserId,
	}

	createdLnkUserProvider, err := s.lnkUserProviderStorage.CreateLnkUserProvider(&createLnkUserProviderPayload, nil)

	if err != nil {
		return 0, fmt.Errorf("%v: create LnkUserProvider error: %w", op, err)
	}

	return createdLnkUserProvider.UserId, nil
}

func (s *AuthLinkService) CreateItem(userId int64) (*models.AuthLink, error) {
	op := "internal.services.authLinksService.CreateItem"

	userProviders, err := s.lnkUserProviderStorage.GetListByUserId(userId)

	if err != nil {
		return nil, fmt.Errorf("%v: %w", op, err)
	}

	if len(*userProviders) != 0 {
		return nil, coreErrors.ErrorUserAlreadyLinked
	}

	randomCode, err := random.GetRandomString(LinkUserCodeLength)

	if err != nil {
		return nil, fmt.Errorf("%v: generate random code error: %w", op, err)
	}

	payload := &models.AuthLinkCreatePayload{
		UserId: userId,
		Code:   randomCode,
		Status: models.AuthLinkStatusCreated,
	}

	authLink, err := s.authLinksStorage.CreateItem(payload)

	if err != nil {
		return nil, fmt.Errorf("%v: create authLink error: %w", op, err)
	}

	return authLink, nil
}
