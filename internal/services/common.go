package services

import (
	"bizarre-vpn-api/internal/storage"
	"bizarre-vpn-api/internal/storage/models"
)

type UserStorage interface {
	GetUserById(ID int64, executor storage.Executor) (*models.BaseUser, error)
	UpdateUserRefreshToken(ID int64, refreshToken string) error
	CreateUser(username string, executor storage.Executor) (userId int64, Err error)
}
