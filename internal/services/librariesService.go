package services

import (
	"bizarre-vpn-api/internal/models"
	"bizarre-vpn-api/internal/shared/logger/sl"
	"fmt"
	"log/slog"
)

type LibraryItemsStorage interface {
	GetList() (*[]models.LibraryItem, error)
	GetItem(id int64) (*models.LibraryItem, error)
}

type LnkProtocolsBackendTypesStorage interface {
	GetProtocolsByBackendTypeId(backendTypeId int64) (*[]models.LibraryItem, error)
}

type LibrariesService struct {
	log                             *slog.Logger
	backendTypesStorage             LibraryItemsStorage
	protocolsStorage                LibraryItemsStorage
	lnkProtocolsBackendTypesStorage LnkProtocolsBackendTypesStorage
}

func NewLibrariesService(
	log *slog.Logger,
	backendTypesStorage LibraryItemsStorage,
	protocolsStorage LibraryItemsStorage,
	lnkProtocolsBackendTypesStorage LnkProtocolsBackendTypesStorage,
) *LibrariesService {
	return &LibrariesService{
		log:                             log,
		backendTypesStorage:             backendTypesStorage,
		protocolsStorage:                protocolsStorage,
		lnkProtocolsBackendTypesStorage: lnkProtocolsBackendTypesStorage,
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

func (service *LibrariesService) GetProtocolsList() (*[]models.LibraryItem, error) {
	op := "internal.services.LibrariesService.GetProtocolsList"

	log := service.log.With(slog.String("op", op))

	list, err := service.protocolsStorage.GetList()
	if err != nil {
		log.Error("error of getting protocols list", sl.Err(err))
		return nil, fmt.Errorf("error of getting protocols list: %w", err)
	}

	return list, nil
}

func (service *LibrariesService) GetProtocolsListByBackendTypeId(backendTypeId int64) (*[]models.LibraryItem, error) {
	op := "internal.services.LibrariesService.GetProtocolsListByBackendTypeId"

	log := service.log.With(slog.String("op", op))

	log.Debug("gettingProtocolsByBackendTypeId", slog.Int64("backendTypeId", backendTypeId))

	_, err := service.backendTypesStorage.GetItem(backendTypeId)

	if err != nil {
		log.Info("backendTypeId is not exist", slog.Int64("backendTypeId", backendTypeId), sl.Err(err))

		return nil, fmt.Errorf("backendTypeId %v is not exist", backendTypeId)
	}

	list, err := service.lnkProtocolsBackendTypesStorage.GetProtocolsByBackendTypeId(backendTypeId)
	if err != nil {
		log.Error("error of getting protocols list", sl.Err(err))
		return nil, fmt.Errorf("error of getting protocols list: %w", err)
	}

	return list, nil
}
