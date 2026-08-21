package repository

import (
	"context"
	"industrial-edge-protocol/internal/domain"
	"sync"
	"time"
)

type Store struct {
	mu       sync.RWMutex
	sites    map[domain.ID]domain.Site
	devices  map[domain.ID]domain.Device
	gateways map[domain.ID]domain.Gateway
	points   map[domain.ID]domain.Point
	plans    map[domain.ID]domain.AcquisitionPlan
	rules    map[domain.ID]domain.Rule
	alarms   map[domain.ID]domain.Alarm
	commands map[domain.ID]domain.Command
	readings []domain.Reading
	rollouts map[domain.ID]domain.Rollout
	next     uint64
}

func NewStore() *Store {
	return &Store{sites: map[domain.ID]domain.Site{}, devices: map[domain.ID]domain.Device{}, gateways: map[domain.ID]domain.Gateway{}, points: map[domain.ID]domain.Point{}, plans: map[domain.ID]domain.AcquisitionPlan{}, rules: map[domain.ID]domain.Rule{}, alarms: map[domain.ID]domain.Alarm{}, commands: map[domain.ID]domain.Command{}, rollouts: map[domain.ID]domain.Rollout{}}
}
func (s *Store) ID(prefix string) domain.ID {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	return domain.ID(prefix + "-" + time.Now().UTC().Format("20060102150405") + "-" + fmtUint(s.next))
}
func fmtUint(v uint64) string {
	const d = "0123456789"
	if v == 0 {
		return "0"
	}
	b := make([]byte, 0, 20)
	for v > 0 {
		b = append([]byte{d[v%10]}, b...)
		v /= 10
	}
	return string(b)
}
func (s *Store) PutSite(_ context.Context, v domain.Site) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sites[v.ID]; ok {
		return domain.ErrConflict
	}
	s.sites[v.ID] = v
	return nil
}
func (s *Store) ListSites(_ context.Context) []domain.Site {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Site, 0, len(s.sites))
	for _, v := range s.sites {
		out = append(out, v)
	}
	return out
}
func (s *Store) PutDevice(_ context.Context, v domain.Device) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.devices[v.ID]; ok {
		return domain.ErrConflict
	}
	s.devices[v.ID] = v
	return nil
}
func (s *Store) ListDevices(_ context.Context) []domain.Device {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Device, 0, len(s.devices))
	for _, v := range s.devices {
		out = append(out, v)
	}
	return out
}
func (s *Store) PutGateway(_ context.Context, v domain.Gateway) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.gateways[v.ID]; ok {
		return domain.ErrConflict
	}
	s.gateways[v.ID] = v
	return nil
}
func (s *Store) ListGateways(_ context.Context) []domain.Gateway {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Gateway, 0, len(s.gateways))
	for _, v := range s.gateways {
		out = append(out, v)
	}
	return out
}
func (s *Store) PutPoint(_ context.Context, v domain.Point) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.points[v.ID]; ok {
		return domain.ErrConflict
	}
	s.points[v.ID] = v
	return nil
}
func (s *Store) ListPoints(_ context.Context) []domain.Point {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Point, 0, len(s.points))
	for _, v := range s.points {
		out = append(out, v)
	}
	return out
}
func (s *Store) PutPlan(_ context.Context, v domain.AcquisitionPlan) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.plans[v.ID]; ok {
		return domain.ErrConflict
	}
	s.plans[v.ID] = v
	return nil
}
func (s *Store) ListPlans(_ context.Context) []domain.AcquisitionPlan {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.AcquisitionPlan, 0, len(s.plans))
	for _, v := range s.plans {
		out = append(out, v)
	}
	return out
}
func (s *Store) PutRule(_ context.Context, v domain.Rule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rules[v.ID]; ok {
		return domain.ErrConflict
	}
	s.rules[v.ID] = v
	return nil
}
func (s *Store) ListRules(_ context.Context) []domain.Rule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Rule, 0, len(s.rules))
	for _, v := range s.rules {
		out = append(out, v)
	}
	return out
}
func (s *Store) PutAlarm(_ context.Context, v domain.Alarm) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.alarms[v.ID] = v
}
func (s *Store) ListAlarms(_ context.Context) []domain.Alarm {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Alarm, 0, len(s.alarms))
	for _, v := range s.alarms {
		out = append(out, v)
	}
	return out
}
func (s *Store) PutCommand(_ context.Context, v domain.Command) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.commands[v.ID]; ok {
		return domain.ErrConflict
	}
	s.commands[v.ID] = v
	return nil
}
func (s *Store) UpdateCommand(_ context.Context, v domain.Command) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.commands[v.ID]
	if !ok {
		return domain.ErrNotFound
	}
	if current.Version != v.Version-1 {
		return domain.ErrConflict
	}
	s.commands[v.ID] = v
	return nil
}
func (s *Store) ListCommands(_ context.Context) []domain.Command {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Command, 0, len(s.commands))
	for _, v := range s.commands {
		out = append(out, v)
	}
	return out
}
func (s *Store) AppendReadings(_ context.Context, values []domain.Reading) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.readings = append(s.readings, values...)
	if len(s.readings) > 10000 {
		s.readings = s.readings[len(s.readings)-10000:]
	}
}
func (s *Store) Readings(_ context.Context, id domain.ID, from, to time.Time) []domain.Reading {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Reading, 0, len(s.readings))
	for _, v := range s.readings {
		if v.PointID == id && !v.ObservedAt.Before(from) && v.ObservedAt.Before(to) {
			out = append(out, v)
		}
	}
	return out
}
func (s *Store) PutRollout(_ context.Context, v domain.Rollout) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rollouts[v.ID] = v
}
func (s *Store) ListRollouts(_ context.Context) []domain.Rollout {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Rollout, 0, len(s.rollouts))
	for _, v := range s.rollouts {
		out = append(out, v)
	}
	return out
}
