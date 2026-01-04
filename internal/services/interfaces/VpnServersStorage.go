package interfaces

import "bizarre-vpn-api/internal/models"

type VpnServersStorage interface {
	GetExpandedList() (*[]models.VpnServerExpandedItem, error)
	GetExpandedItemById(vpnServerId int64) (*models.VpnServerExpandedItem, error)
	CreateItem(vpnServer *models.VpnServerItem) (int64, error)
	DeleteItem(vpnServerId int64) error
}
