package models

type LnkProtocolsBackendTypes struct {
	ID            int64 `json:"id" db:"id"`
	ProtocolId    int64 `json:"protocolId" db:"protocol_id"`
	BackendTypeId int64 `json:"backendTypeId" db:"backend_type_id"`
}
