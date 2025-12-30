package models

type AuthLinkStatus = string

const (
	AuthLinkStatusCreated = "created"
)

type AuthLink struct {
	Id     int64          `json:"id" db:"id"`
	UserId int64          `json:"userId" db:"user_id"`
	Code   string         `json:"code" db:"code"`
	Status AuthLinkStatus `json:"status" db:"status"`
}
