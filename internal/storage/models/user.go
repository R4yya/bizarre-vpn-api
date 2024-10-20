package models

// User represents the Telegram user
type User struct {
	ID           int64  `db:"id"`
	TelegramID   int64  `db:"telegram_id"`
	Username     string `db:"username"`
	LanguageCode string `db:"language_code"`
	IsBot        bool   `db:"is_bot"`
}
