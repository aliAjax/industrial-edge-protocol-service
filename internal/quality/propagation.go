package quality

import "industrial-edge-protocol/internal/domain"

func Worst(values ...domain.Quality) domain.Quality {
	result := domain.QualityGood
	rank := map[domain.Quality]int{domain.QualityGood: 0, domain.QualityUncertain: 1, domain.QualityStale: 2, domain.QualityBad: 3}
	for _, v := range values {
		if rank[v] > rank[result] {
			result = v
		}
	}
	return result
}
func Transform(value domain.Reading, gain, offset float64) domain.Reading {
	value.Value = value.Value*gain + offset
	if value.Quality == "" {
		value.Quality = domain.QualityUncertain
	}
	return value
}
func Stale(value domain.Reading, maxAge, now int64) domain.Reading {
	if value.ObservedAt.Unix()+maxAge < now {
		value.Quality = domain.QualityStale
	}
	return value
}
