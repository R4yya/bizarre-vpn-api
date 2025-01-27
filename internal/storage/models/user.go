package models

import "time"

// User represents the Telegram user
type BaseUser struct {
	ID        int64     `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Role      string    `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
type FullUser struct {
	BaseUser
	RefreshToken string `db:"refresh_token" json:"refreshToken"`
}
