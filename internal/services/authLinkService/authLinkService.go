package authLinkService

import (
	"bizarre-vpn-api/internal/lib/random"
	"bizarre-vpn-api/internal/services"
	"bizarre-vpn-api/internal/storage/models"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrUserAlreadyLinked = errors.New("user already linked")
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

func (s *AuthLinkService) GetItemByCode(code string) (*models.AuthLink, error) {
	return s.authLinksStorage.GetItemByCode(code)
}

func (s *AuthLinkService) CreateItem(userId int64) (*models.AuthLink, error) {
	op := "internal.services.authLinksService.CreateItem"

	userProviders, err := s.lnkUserProviderStorage.GetListByUserId(userId)

	if err != nil {
		return nil, fmt.Errorf("%v: %w", op, err)
	}

	if len(*userProviders) != 0 {
		return nil, ErrUserAlreadyLinked
	}

	randomCode, err := random.GetRandomString(6)

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
