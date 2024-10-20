package models

// User represents the Telegram user
type User struct {
	ID           int64  `db:"id" json:"id"`
	TelegramID   int64  `db:"telegram_id" json:"telegramID"`
	Username     string `db:"username" json:"username"`
	LanguageCode string `db:"language_code" json:"languageCode"`
	IsBot        bool   `db:"is_bot" json:"isBot"`
}
