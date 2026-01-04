package interfaces

import (
	"bizarre-vpn-api/internal/models"
	"bizarre-vpn-api/internal/storage"
)

type LnkUserProviderStorage interface {
	GetItemByType(
		providerType string,
		externalId string,
	) (LnkUserProvider *models.LnkUserProvider, isFound bool, Err error)

	GetListByUserId(userId int64) (*[]models.LnkUserProvider, error)

	CreateLnkUserProvider(
		createLnkUserProviderPayload *models.CreateLnkUserProviderPayload,
		executor storage.Executor,
	) (*models.LnkUserProvider, error)
}
