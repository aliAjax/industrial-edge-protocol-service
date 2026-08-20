package replay

import (
	"context"
	"industrial-edge-protocol/internal/domain"
	"sort"
	"time"
)

type Handler func(context.Context, domain.Reading) error

func Run(ctx context.Context, values []domain.Reading, handler Handler, pace time.Duration) error {
	sorted := append([]domain.Reading(nil), values...)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].ObservedAt.Before(sorted[j].ObservedAt) })
	var previous time.Time
	for _, value := range sorted {
		if !previous.IsZero() && pace > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(pace):
			}
		}
		if err := handler(ctx, value); err != nil {
			return err
		}
		previous = value.ObservedAt
	}
	return nil
}
func GroupByPoint(values []domain.Reading) map[domain.ID][]domain.Reading {
	out := map[domain.ID][]domain.Reading{}
	for _, v := range values {
		out[v.PointID] = append(out[v.PointID], v)
	}
	return out
}
