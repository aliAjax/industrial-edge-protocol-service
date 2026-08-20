package quality

import (
	"industrial-edge-protocol/internal/domain"
	"math"
)

type Range struct {
	Min float64
	Max float64
}
type Validator struct {
	Ranges map[domain.ID]Range
	Last   map[domain.ID]domain.Reading
}

func NewValidator() *Validator {
	return &Validator{Ranges: map[domain.ID]Range{}, Last: map[domain.ID]domain.Reading{}}
}
func (v *Validator) Validate(r domain.Reading) domain.Reading {
	if math.IsNaN(r.Value) || math.IsInf(r.Value, 0) {
		r.Quality = domain.QualityBad
		return r
	}
	if bounds, ok := v.Ranges[r.PointID]; ok && (r.Value < bounds.Min || r.Value > bounds.Max) {
		r.Quality = domain.QualityUncertain
	}
	if prev, ok := v.Last[r.PointID]; ok && r.ObservedAt.Before(prev.ObservedAt) {
		r.Quality = domain.QualityUncertain
	}
	v.Last[r.PointID] = r
	return r
}
func (v *Validator) SetRange(id domain.ID, r Range) { v.Ranges[id] = r }
