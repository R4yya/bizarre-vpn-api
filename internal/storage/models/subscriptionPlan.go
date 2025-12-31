package models

import "time"

// SubscriptionPlan presents a subscription plan
type SubscriptionPlan struct {
	ID             int64     `db:"id" json:"id"`
	Country        string    `db:"country" json:"country"`
	Name           string    `db:"name" json:"name"`
	Description    string    `db:"description" json:"description"`
	DurationMonths int       `db:"duration_months" json:"durationMonths"`
	DataLimitGB    *int      `db:"data_limit_gb" json:"dataLimitGb"`       // NULL for unlimited use
	SpeedLimitMbps *int      `db:"speed_limit_mbps" json:"speedLimitMbps"` // NULL for unlimited
	DeviceLimit    int       `db:"device_limit" json:"deviceLimit"`
	Price          float64   `db:"price" json:"price"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time `db:"updated_at" json:"updatedAt"`
}
