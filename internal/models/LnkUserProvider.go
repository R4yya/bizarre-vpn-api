package models

type Provider = string

const (
	TelegramProviderName Provider = "telegram"
)

type LnkUserProvider struct {
	ID             int64  `json:"id" db:"id"`
	ProviderType   string `json:"providerType" db:"provider_type"`
	ExternalUserId string `json:"externalUserId" db:"external_user_id"`
	UserId         int64  `json:"userId" db:"user_id"`
}

type CreateLnkUserProviderPayload struct {
	ProviderType   string `json:"providerType" db:"provider_type"`
	ExternalUserId string `json:"externalUserId" db:"external_user_id"`
	UserId         int64  `json:"userId" db:"user_id"`
}
