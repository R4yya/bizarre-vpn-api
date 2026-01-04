package interfaces

import (
	"bizarre-vpn-api/internal/models"
	"bizarre-vpn-api/internal/storage"
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
