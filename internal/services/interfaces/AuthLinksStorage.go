package interfaces

import "bizarre-vpn-api/internal/models"

type AuthLinksStorage interface {
	GetItemByCode(code string) (*models.AuthLink, error)
	CreateItem(payload *models.AuthLinkCreatePayload) (*models.AuthLink, error)
}
