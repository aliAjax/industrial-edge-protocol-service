package api

import "time"

type ErrorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}
type CreatePointRequest struct {
	DeviceID    string        `json:"device_id"`
	Name        string        `json:"name"`
	Address     uint16        `json:"address"`
	DataType    string        `json:"data_type"`
	Unit        string        `json:"unit"`
	Scale       float64       `json:"scale"`
	Deadband    float64       `json:"deadband"`
	SampleEvery time.Duration `json:"sample_every"`
}
type IngestRequest struct {
	PointID    string    `json:"point_id"`
	Value      float64   `json:"value"`
	Quality    string    `json:"quality"`
	ObservedAt time.Time `json:"observed_at"`
}
type CommandRequest struct {
	DeviceID    string  `json:"device_id"`
	PointID     string  `json:"point_id"`
	Value       float64 `json:"value"`
	DryRun      bool    `json:"dry_run"`
	Interlock   string  `json:"interlock"`
	RequestedBy string  `json:"requested_by"`
}
type RolloutRequest struct {
	ConfigVersion int64    `json:"config_version"`
	GatewayIDs    []string `json:"gateway_ids"`
	Percent       int      `json:"percent"`
}
