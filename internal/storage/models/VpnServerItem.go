package models

import "time"

type VpnServerItem struct {
	ID               int64     `db:"id" json:"id"`
	Name             string    `db:"name" json:"name"`
	Description      string    `db:"description" json:"description"`
	AdapterHost      string    `db:"adapter_host" json:"adapterHost"`
	AdapterPort      uint16    `db:"adapter_port" json:"adapterPort"`
	Country          string    `db:"country" json:"country"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time `db:"updated_at" json:"updatedAt"`
	LnkBackendTypeId int64     `db:"lnk_backend_type_id" json:"lnkBackendTypeId"`
}

type VpnServerExpandedItem struct {
	VpnServerItem
	Protocol    string `db:"protocol" json:"protocol"`
	BackendType string `db:"backend_type" json:"backendType"`
}
