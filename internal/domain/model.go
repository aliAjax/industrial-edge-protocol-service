package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrNotFound      = errors.New("entity not found")
	ErrConflict      = errors.New("version conflict")
	ErrInvalid       = errors.New("invalid entity")
	ErrUnsafeCommand = errors.New("unsafe command")
)

type ID string
type Quality string

const (
	QualityGood      Quality = "good"
	QualityBad       Quality = "bad"
	QualityUncertain Quality = "uncertain"
	QualityStale     Quality = "stale"
)

type Site struct {
	ID        ID        `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	Version   int64     `json:"version"`
}
type Device struct {
	ID       ID     `json:"id"`
	SiteID   ID     `json:"site_id"`
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
	Enabled  bool   `json:"enabled"`
	Version  int64  `json:"version"`
}
type Gateway struct {
	ID           ID        `json:"id"`
	SiteID       ID        `json:"site_id"`
	Name         string    `json:"name"`
	Capabilities []string  `json:"capabilities"`
	Connected    bool      `json:"connected"`
	LastSeen     time.Time `json:"last_seen"`
	Version      int64     `json:"version"`
}
type Point struct {
	ID          ID            `json:"id"`
	DeviceID    ID            `json:"device_id"`
	Name        string        `json:"name"`
	Address     uint16        `json:"address"`
	DataType    string        `json:"data_type"`
	Unit        string        `json:"unit"`
	Scale       float64       `json:"scale"`
	Deadband    float64       `json:"deadband"`
	SampleEvery time.Duration `json:"sample_every"`
	Enabled     bool          `json:"enabled"`
	Version     int64         `json:"version"`
}
type Calibration struct {
	ID          ID        `json:"id"`
	PointID     ID        `json:"point_id"`
	Offset      float64   `json:"offset"`
	Gain        float64   `json:"gain"`
	EffectiveAt time.Time `json:"effective_at"`
	Version     int64     `json:"version"`
}
type Reading struct {
	PointID    ID        `json:"point_id"`
	Value      float64   `json:"value"`
	Unit       string    `json:"unit"`
	Quality    Quality   `json:"quality"`
	ObservedAt time.Time `json:"observed_at"`
	ReceivedAt time.Time `json:"received_at"`
	Sequence   uint64    `json:"sequence"`
}
type AcquisitionPlan struct {
	ID         ID            `json:"id"`
	GatewayID  ID            `json:"gateway_id"`
	PointIDs   []ID          `json:"point_ids"`
	Interval   time.Duration `json:"interval"`
	Timeout    time.Duration `json:"timeout"`
	RetryLimit int           `json:"retry_limit"`
	Enabled    bool          `json:"enabled"`
	Version    int64         `json:"version"`
}
type RuleState string

const (
	RuleNormal     RuleState = "normal"
	RulePending    RuleState = "pending"
	RuleAlarm      RuleState = "alarm"
	RuleSuppressed RuleState = "suppressed"
	RuleCleared    RuleState = "cleared"
)

type Rule struct {
	ID          ID            `json:"id"`
	Name        string        `json:"name"`
	Expression  string        `json:"expression"`
	Window      int           `json:"window"`
	Threshold   float64       `json:"threshold"`
	Hysteresis  float64       `json:"hysteresis"`
	Debounce    time.Duration `json:"debounce"`
	Cooldown    time.Duration `json:"cooldown"`
	Maintenance bool          `json:"maintenance"`
	Enabled     bool          `json:"enabled"`
	Version     int64         `json:"version"`
}
type Alarm struct {
	ID          ID         `json:"id"`
	RuleID      ID         `json:"rule_id"`
	PointID     ID         `json:"point_id"`
	State       RuleState  `json:"state"`
	TriggeredAt time.Time  `json:"triggered_at"`
	ClearedAt   *time.Time `json:"cleared_at,omitempty"`
	LastValue   float64    `json:"last_value"`
	Count       int        `json:"count"`
}
type CommandStatus string

const (
	CommandDraft     CommandStatus = "draft"
	CommandPending   CommandStatus = "pending"
	CommandApproved  CommandStatus = "approved"
	CommandSent      CommandStatus = "sent"
	CommandConfirmed CommandStatus = "confirmed"
	CommandRejected  CommandStatus = "rejected"
	CommandExpired   CommandStatus = "expired"
)

type Command struct {
	ID          ID            `json:"id"`
	DeviceID    ID            `json:"device_id"`
	PointID     ID            `json:"point_id"`
	Value       float64       `json:"value"`
	DryRun      bool          `json:"dry_run"`
	Interlock   string        `json:"interlock"`
	Status      CommandStatus `json:"status"`
	RequestedBy string        `json:"requested_by"`
	ApprovedBy  string        `json:"approved_by"`
	ExpiresAt   time.Time     `json:"expires_at"`
	CreatedAt   time.Time     `json:"created_at"`
	Version     int64         `json:"version"`
}
type Rollout struct {
	ID            ID        `json:"id"`
	ConfigVersion int64     `json:"config_version"`
	GatewayIDs    []ID      `json:"gateway_ids"`
	Percent       int       `json:"percent"`
	State         string    `json:"state"`
	Error         string    `json:"error,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
type TelemetryAggregate struct {
	PointID  ID        `json:"point_id"`
	From     time.Time `json:"from"`
	To       time.Time `json:"to"`
	Count    int       `json:"count"`
	Min      float64   `json:"min"`
	Max      float64   `json:"max"`
	Average  float64   `json:"average"`
	BadCount int       `json:"bad_count"`
}

func ValidateIdentifier(value string) error {
	value = strings.TrimSpace(value)
	if value == "" || len(value) > 128 {
		return fmt.Errorf("%w: identifier", ErrInvalid)
	}
	return nil
}
func (p Point) Validate() error {
	if p.DeviceID == "" || p.Name == "" || p.SampleEvery <= 0 || p.Scale == 0 {
		return ErrInvalid
	}
	return nil
}
func (r Rule) Validate() error {
	if r.ID == "" || r.Expression == "" || r.Window <= 0 || r.Hysteresis < 0 {
		return ErrInvalid
	}
	return nil
}
func (c Command) Validate() error {
	if c.DeviceID == "" || c.PointID == "" || c.ExpiresAt.Before(time.Now()) {
		return ErrInvalid
	}
	if c.Value != c.Value {
		return ErrInvalid
	}
	return nil
}
