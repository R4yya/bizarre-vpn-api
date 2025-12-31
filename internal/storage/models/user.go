package models

import "time"

type UserRole = string

const (
	UserRoleClient = "client"
	UserRoleAdmin  = "admin"
)

// User represents the Telegram user
type BaseUser struct {
	ID        int64     `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Login     *string   `db:"login" json:"login"`
	Role      UserRole  `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
type FullUser struct {
	BaseUser
	RefreshToken string `db:"refresh_token" json:"refreshToken"`
}

type CreateUserPayload struct {
	Login    *string  `json:"login" validate:"required"`
	Username string   `json:"username" validate:"required"`
	Role     UserRole `json:"role" validate:"required"`
	Password *string  `json:"password" validate:"required"`
}
