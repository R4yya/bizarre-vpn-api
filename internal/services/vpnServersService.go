package services

import (
	"bizarre-vpn-api/internal/models"
	"bizarre-vpn-api/internal/services/interfaces"
	"bizarre-vpn-api/internal/shared/coreErrors"
	"fmt"
	"log/slog"
	"net"
	"unicode/utf8"
)

type VpnServersService struct {
	log               *slog.Logger
	vpnServersStorage interfaces.VpnServersStorage
}

func NewVpnServersService(
	log *slog.Logger,
	vpnServersStorage interfaces.VpnServersStorage,
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
		return 0, coreErrors.ValidationError{
			Msg:    "invalid IP",
			Entity: "vpnServer",
			Fields: []string{
				"adapterHost",
			},
		}
	}

	if vpnServer.AdapterPort == 0 {
		return 0, coreErrors.ValidationError{
			Msg:    "invalid port",
			Entity: "vpnServer",
			Fields: []string{
				"adapterPort",
			},
		}
	}

	if nameLen := utf8.RuneCountInString(vpnServer.Name); nameLen < 3 || nameLen >= 250 {
		return 0, coreErrors.ValidationError{
			Msg:    "must be grater than 3 and least than 250",
			Entity: "vpnServer",
			Fields: []string{
				"name",
			},
		}
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
