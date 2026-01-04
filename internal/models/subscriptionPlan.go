package models

import (
	"time"
)

type VpnServerID = int64

type SubscriptionPlan struct {
	ID           int64       `db:"id" json:"id"`
	Name         string      `db:"name" json:"name"`
	Description  string      `db:"description" json:"description"`
	Price        int64       `db:"price" json:"price"`
	DurationDays int32       `db:"duration_days" json:"durationDays"`
	VpnServerID  VpnServerID `db:"vpn_server_id" json:"vpnServerId"`
	CreatedAt    time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time   `db:"updated_at" json:"updatedAt"`
}

type CreateSubscriptionPlan struct {
	Name         string `db:"name" json:"name"`
	Description  string `db:"description" json:"description"`
	Price        int64  `db:"price" json:"price"`
	DurationDays int32  `db:"duration_days" json:"durationDays"`
	VpnServerID  int64  `db:"vpn_server_id" json:"vpnServerId"`
}

type UpdateSubscriptionPlan struct {
	Name         string `db:"name" json:"name"`
	Description  string `db:"description" json:"description"`
	Price        int64  `db:"price" json:"price"`
	DurationDays int32  `db:"duration_days" json:"durationDays"`
}
