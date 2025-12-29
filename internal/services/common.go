package services

import (
	"bizarre-vpn-api/internal/storage"
	"bizarre-vpn-api/internal/storage/models"
)

type UserStorage interface {
	GetUsersList() (*[]models.BaseUser, error)
	GetUserById(ID int64, executor storage.Executor) (*models.BaseUser, error)
	GetUserByUsername(userName string) (*models.BaseUser, error)
	GetUserPasswordHash(userId int64) (string, error)
	GetUserRefreshToken(ID int64) (string, error)
	CreateUser(payload *models.CreateUserPayload, executor storage.Executor) (*models.BaseUser, error)
	UpdateUserRefreshToken(ID int64, refreshToken string) error
}
