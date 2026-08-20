package telemetry

import (
	"industrial-edge-protocol/internal/domain"
	"sort"
	"sync"
	"time"
)

type Store struct {
	mu      sync.RWMutex
	byPoint map[domain.ID][]domain.Reading
}

func NewStore() *Store { return &Store{byPoint: map[domain.ID][]domain.Reading{}} }
func (s *Store) Append(values []domain.Reading) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range values {
		s.byPoint[v.PointID] = append(s.byPoint[v.PointID], v)
		sort.SliceStable(s.byPoint[v.PointID], func(i, j int) bool {
			return s.byPoint[v.PointID][i].ObservedAt.Before(s.byPoint[v.PointID][j].ObservedAt)
		})
	}
}
func (s *Store) Query(id domain.ID, from, to time.Time, limit int) []domain.Reading {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.Reading{}
	for _, v := range s.byPoint[id] {
		if !v.ObservedAt.Before(from) && v.ObservedAt.Before(to) {
			out = append(out, v)
			if limit > 0 && len(out) >= limit {
				break
			}
		}
	}
	return out
}
func (s *Store) Aggregate(id domain.ID, from, to time.Time) domain.TelemetryAggregate {
	rows := s.Query(id, from, to, 0)
	a := domain.TelemetryAggregate{PointID: id, From: from, To: to, Count: len(rows)}
	if len(rows) == 0 {
		return a
	}
	a.Min = rows[0].Value
	a.Max = rows[0].Value
	for _, v := range rows {
		if v.Value < a.Min {
			a.Min = v.Value
		}
		if v.Value > a.Max {
			a.Max = v.Value
		}
		a.Average += v.Value
		if v.Quality != domain.QualityGood {
			a.BadCount++
		}
	}
	a.Average /= float64(len(rows))
	return a
}
