package acquisition

import (
	"context"
	"industrial-edge-protocol/internal/domain"
	"math/rand"
	"time"
)

type MockReader struct {
	Seed      int64
	FailEvery int
}

func (m MockReader) Read(ctx context.Context, p domain.Point) (domain.Reading, error) {
	v := rand.New(rand.NewSource(m.Seed+int64(p.Address)+time.Now().Unix()/60)).Float64() * 100
	return domain.Reading{PointID: p.ID, Value: v * p.Scale, Unit: p.Unit, Quality: domain.QualityGood, ObservedAt: time.Now().UTC()}, nil
}
