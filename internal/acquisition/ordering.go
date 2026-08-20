package acquisition

import (
	"industrial-edge-protocol/internal/domain"
	"sort"
)

func Order(values []domain.Reading) []domain.Reading {
	out := append([]domain.Reading(nil), values...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].ObservedAt.Equal(out[j].ObservedAt) {
			return out[i].Sequence < out[j].Sequence
		}
		return out[i].ObservedAt.Before(out[j].ObservedAt)
	})
	return out
}
func Deduplicate(values []domain.Reading) []domain.Reading {
	out := []domain.Reading{}
	seen := map[uint64]bool{}
	for _, v := range Order(values) {
		if v.Sequence > 0 && seen[v.Sequence] {
			continue
		}
		if v.Sequence > 0 {
			seen[v.Sequence] = true
		}
		out = append(out, v)
	}
	return out
}
