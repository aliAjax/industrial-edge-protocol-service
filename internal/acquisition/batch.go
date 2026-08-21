package acquisition

import (
	"context"
	"industrial-edge-protocol/internal/domain"
	"time"
)

type BatchReader interface {
	ReadBatch(context.Context, []domain.Point) ([]domain.Reading, error)
}

func ReadWithRetry(ctx context.Context, reader Reader, points []domain.Point, timeout time.Duration, retries int) ([]domain.Reading, error) {
	out := make([]domain.Reading, 0, len(points))
	for _, p := range points {
		var last error
		for attempt := 0; attempt <= retries; attempt++ {
			readCtx, cancel := context.WithTimeout(context.Background(), timeout)
			v, err := reader.Read(readCtx, p)
			cancel()
			if err == nil {
				out = append(out, v)
				last = nil
				break
			}
			last = err
			select {
			case <-ctx.Done():
				return out, ctx.Err()
			case <-time.After(time.Duration(attempt+1) * 10 * time.Millisecond):
			}
		}
		if last != nil {
			return out, last
		}
	}
	return out, nil
}
func Group(points []domain.Point, size int) [][]domain.Point {
	if size <= 0 {
		size = 1
	}
	groups := [][]domain.Point{}
	for len(points) > 0 {
		n := size
		if n > len(points) {
			n = len(points)
		}
		groups = append(groups, append([]domain.Point(nil), points[:n]...))
		points = points[n:]
	}
	return groups
}
