package service

import (
	"context"
	"fmt"
	"industrial-edge-protocol/internal/buffer"
	"industrial-edge-protocol/internal/config"
	"industrial-edge-protocol/internal/domain"
	"industrial-edge-protocol/internal/repository"
	"industrial-edge-protocol/internal/rules"
	"industrial-edge-protocol/internal/telemetry"
	"log/slog"
	"strings"
	"sync"
	"time"
)

type Application struct {
	Log       *slog.Logger
	Config    config.Config
	Repo      *repository.Store
	WAL       *buffer.WAL
	Telemetry *telemetry.Store
	Rules     *rules.Engine
	mu        sync.Mutex
	sequence  uint64
}

func NewApplication(log *slog.Logger, cfg config.Config) *Application {
	repo := repository.NewStore()
	return &Application{Log: log, Config: cfg, Repo: repo, WAL: buffer.New(cfg.DataDir, cfg.MaxBufferBytes), Telemetry: telemetry.NewStore(), Rules: rules.NewEngine()}
}
func (a *Application) CreateSite(ctx context.Context, v domain.Site) (domain.Site, error) {
	if err := domain.ValidateIdentifier(v.Name); err != nil {
		return v, err
	}
	v.ID = a.Repo.ID("site")
	v.CreatedAt = time.Now().UTC()
	v.Version = 1
	return v, a.Repo.PutSite(ctx, v)
}
func (a *Application) CreateDevice(ctx context.Context, v domain.Device) (domain.Device, error) {
	if err := domain.ValidateIdentifier(v.Name); err != nil {
		return v, err
	}
	if v.SiteID == "" {
		return v, domain.ErrInvalid
	}
	v.ID = a.Repo.ID("device")
	v.Version = 1
	return v, a.Repo.PutDevice(ctx, v)
}
func (a *Application) CreateGateway(ctx context.Context, v domain.Gateway) (domain.Gateway, error) {
	if err := domain.ValidateIdentifier(v.Name); err != nil {
		return v, err
	}
	v.ID = a.Repo.ID("gateway")
	v.LastSeen = time.Now().UTC()
	v.Version = 1
	return v, a.Repo.PutGateway(ctx, v)
}
func (a *Application) CreatePoint(ctx context.Context, v domain.Point) (domain.Point, error) {
	if err := v.Validate(); err != nil {
		return v, err
	}
	v.ID = a.Repo.ID("point")
	v.Version = 1
	return v, a.Repo.PutPoint(ctx, v)
}
func (a *Application) CreatePlan(ctx context.Context, v domain.AcquisitionPlan) (domain.AcquisitionPlan, error) {
	if v.GatewayID == "" || len(v.PointIDs) == 0 || v.Interval <= 0 {
		return v, domain.ErrInvalid
	}
	v.ID = a.Repo.ID("plan")
	v.Enabled = true
	v.Version = 1
	return v, a.Repo.PutPlan(ctx, v)
}
func (a *Application) CreateRule(ctx context.Context, v domain.Rule) (domain.Rule, error) {
	if err := v.Validate(); err != nil {
		return v, err
	}
	v.ID = a.Repo.ID("rule")
	v.Enabled = true
	v.Version = 1
	return v, a.Repo.PutRule(ctx, v)
}
func (a *Application) CreateCommand(ctx context.Context, v domain.Command) (domain.Command, error) {
	if v.ExpiresAt.IsZero() {
		v.ExpiresAt = time.Now().Add(a.Config.CommandTimeout)
	}
	v.Status = domain.CommandDraft
	v.CreatedAt = time.Now().UTC()
	v.ID = a.Repo.ID("command")
	v.Version = 1
	if err := v.Validate(); err != nil {
		return v, err
	}
	return v, a.Repo.PutCommand(ctx, v)
}
func (a *Application) ApproveCommand(ctx context.Context, id, actor string) (domain.Command, error) {
	for _, c := range a.Repo.ListCommands(ctx) {
		if c.ID == domain.ID(id) {
			if c.Status != domain.CommandDraft && c.Status != domain.CommandPending {
				return c, domain.ErrConflict
			}
			if strings.TrimSpace(actor) == "" {
				return c, domain.ErrInvalid
			}
			c.Status = domain.CommandApproved
			c.ApprovedBy = actor
			c.Version++
			return c, a.Repo.UpdateCommand(ctx, c)
		}
	}
	return domain.Command{}, domain.ErrNotFound
}
func (a *Application) Ingest(ctx context.Context, values []domain.Reading) error {
	if len(values) == 0 {
		return nil
	}
	a.mu.Lock()
	for i := range values {
		a.sequence++
		values[i].Sequence = a.sequence
		if values[i].ReceivedAt.IsZero() {
			values[i].ReceivedAt = time.Now().UTC()
		}
	}
	a.mu.Unlock()
	if _, err := a.WAL.Append(values); err != nil {
		return err
	}
	a.Telemetry.Append(values)
	a.Repo.AppendReadings(ctx, values)
	for _, r := range values {
		for _, rule := range a.Repo.ListRules(ctx) {
			alarm, ok, err := a.Rules.Evaluate(rule, r)
			if err != nil {
				a.Log.Warn("rule evaluation failed", "error", err)
				continue
			}
			if ok {
				a.Repo.PutAlarm(ctx, alarm)
			}
		}
	}
	return nil
}
func (a *Application) RunCycle(ctx context.Context) error {
	for _, p := range a.Repo.ListPlans(ctx) {
		if !p.Enabled {
			continue
		}
		values := make([]domain.Reading, 0, len(p.PointIDs))
		for _, id := range p.PointIDs {
			values = append(values, domain.Reading{PointID: id, Value: float64(time.Now().UnixNano()%1000) / 10, Quality: domain.QualityGood, ObservedAt: time.Now().UTC()})
		}
		if err := a.Ingest(ctx, values); err != nil {
			return fmt.Errorf("plan %s: %w", p.ID, err)
		}
	}
	return nil
}
func (a *Application) Health() map[string]any {
	return map[string]any{"status": "ok", "time": time.Now().UTC(), "buffer_dir": a.Config.DataDir, "points": len(a.Repo.ListPoints(context.Background()))}
}
