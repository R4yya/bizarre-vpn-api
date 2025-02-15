package services

import (
	"bizarre-vpn-api/internal/lib/logger/sl"
	"bizarre-vpn-api/internal/storage/models"
	"fmt"
	"log/slog"
)

type BackendTypesStorage interface {
	GetList() (*[]models.LibraryItem, error)
}

type LibrariesService struct {
	log                 *slog.Logger
	backendTypesStorage BackendTypesStorage
}

func NewLibrariesService(
	log *slog.Logger,
	backendTypesStorage BackendTypesStorage,
) *LibrariesService {
	return &LibrariesService{
		log:                 log,
		backendTypesStorage: backendTypesStorage,
	}
}

func (service *LibrariesService) GetBackendTypesList() (*[]models.LibraryItem, error) {
	op := "internal.services.LibrariesService.GetBackendTypesList"

	log := service.log.With(slog.String("op", op))

	list, err := service.backendTypesStorage.GetList()
	if err != nil {
		log.Error("error of getting backendTypes list", sl.Err(err))
		return nil, fmt.Errorf("error of getting backendTypes list: %w", err)
	}

	return list, nil
}
