package services

import (
	"bizarre-vpn-api/internal/core/coreErrors"
	"bizarre-vpn-api/internal/storage/models"
	"fmt"
	"log/slog"
	"net"
	"unicode/utf8"
)

type VpnServersStorage interface {
	GetExpandedList() (*[]models.VpnServerExpandedItem, error)
	GetExpandedItemById(vpnServerId int64) (*models.VpnServerExpandedItem, error)
	CreateItem(vpnServer *models.VpnServerItem) (int64, error)
	DeleteItem(vpnServerId int64) error
}

type VpnServersService struct {
	log               *slog.Logger
	vpnServersStorage VpnServersStorage
}

func NewVpnServersService(
	log *slog.Logger,
	vpnServersStorage VpnServersStorage,
) *VpnServersService {
	return &VpnServersService{
		log:               log,
		vpnServersStorage: vpnServersStorage,
	}
}

func (service *VpnServersService) GetExpandedList() (*[]models.VpnServerExpandedItem, error) {
	op := "internal.services.vpnServersStorage.GetExpandedList"

	list, err := service.vpnServersStorage.GetExpandedList()

	if err != nil {
		return nil, fmt.Errorf("%v: %w", op, err)
	}

	return list, nil
}

func (service *VpnServersService) GetExpandedItemById(vpnServerId int64) (*models.VpnServerExpandedItem, error) {
	op := "internal.services.vpnServersStorage.GetExpandedItemById"

	item, err := service.vpnServersStorage.GetExpandedItemById(vpnServerId)

	if err != nil {
		return nil, fmt.Errorf("%v: %w", op, err)
	}

	return item, nil
}

func (service *VpnServersService) CreateItem(vpnServer *models.VpnServerItem) (int64, error) {
	op := "internal.services.vpnServersStorage.CreateItem"

	if parsedIp := net.ParseIP(vpnServer.AdapterHost); parsedIp == nil {
		return 0, coreErrors.ErrorVpnServerInvalidIp
	}

	if vpnServer.AdapterPort == 0 {
		return 0, coreErrors.ErrorVpnServerInvalidPort
	}

	if nameLen := utf8.RuneCountInString(vpnServer.Name); nameLen < 3 {
		return 0, coreErrors.ErrorVpnServerNameIsTooSmall
	}

	vpnServerID, err := service.vpnServersStorage.CreateItem(vpnServer)

	if err != nil {
		return 0, fmt.Errorf("%v: %w", op, err)
	}

	return vpnServerID, nil
}

func (service *VpnServersService) DeleteItem(vpnServerId int64) error {
	//op := "internal.services.vpnServersStorage.DeleteItem"

	err := service.vpnServersStorage.DeleteItem(vpnServerId)

	if err != nil {
		return err
	}

	return nil
}
