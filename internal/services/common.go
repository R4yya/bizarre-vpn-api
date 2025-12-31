package services

import (
	"bizarre-vpn-api/internal/storage"
	"bizarre-vpn-api/internal/storage/models"
)

type UserStorage interface {
	GetUsersList() (*[]models.BaseUser, error)
	GetUserById(ID int64, executor storage.Executor) (*models.BaseUser, error)
	GetUserByLogin(login string) (*models.BaseUser, error)
	GetUserPasswordHash(userId int64) (string, error)
	GetUserRefreshToken(ID int64) (string, error)
	CreateUser(payload *models.CreateUserPayload, executor storage.Executor) (*models.BaseUser, error)
	UpdateUserRefreshToken(ID int64, refreshToken string) error
}

type LnkUserProviderStorage interface {
	GetItemByType(
		providerType string,
		externalId string,
	) (LnkUserProvider *models.LnkUserProvider, isFound bool, Err error)

	GetListByUserId(userId int64) (*[]models.LnkUserProvider, error)

	CreateLnkUserProvider(
		lnkUserProvider *models.LnkUserProvider,
		username string,
	) (providerId int64, createdUser *models.BaseUser, Err error)
}
