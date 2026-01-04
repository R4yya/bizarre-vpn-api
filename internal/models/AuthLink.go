package models

import "time"

type AuthLinkStatus = string

const (
	AuthLinkStatusCreated = "created"
)

type AuthLink struct {
	Id        int64          `db:"id" json:"id" `
	UserId    int64          `db:"user_id" json:"userId" `
	Code      string         `db:"code" json:"code" `
	Status    AuthLinkStatus `db:"status" json:"status" `
	CreatedAt time.Time      `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time      `db:"updated_at" json:"updatedAt"`
}

type AuthLinkCreatePayload struct {
	UserId int64          `db:"user_id" json:"userId"`
	Code   string         `db:"code" json:"code"`
	Status AuthLinkStatus `db:"status" json:"status" `
}
