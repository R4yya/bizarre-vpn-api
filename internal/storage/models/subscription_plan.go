package models

// SubscriptionPlan presents a subscription plan
type SubscriptionPlan struct {
	ID             int64   `db:"id" json:"id"`
	Country        string  `db:"country" json:"country"`
	Name           string  `db:"name" json:"name"`
	Description    string  `db:"description" json:"description"`
	DurationMonths int     `db:"duration_months" json:"duration_months"`
	DataLimitGB    *int    `db:"data_limit_gb" json:"data_limit_gb"`       // NULL for unlimited use
	SpeedLimitMbps *int    `db:"speed_limit_mbps" json:"speed_limit_mbps"` // NULL for unlimited
	DeviceLimit    int     `db:"device_limit" json:"device_limit"`
	Price          float64 `db:"price" json:"price"`
}
