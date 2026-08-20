package aggregation

import "math"

type Detector struct {
	Mean      float64
	StdDev    float64
	Threshold float64
	Count     int
}

func (d *Detector) Observe(value float64) bool {
	d.Count++
	if d.Count == 1 {
		d.Mean = value
		return false
	}
	delta := value - d.Mean
	d.Mean += delta / float64(d.Count)
	d.StdDev += delta * (value - d.Mean)
	if d.Count < 4 {
		return false
	}
	variance := d.StdDev / float64(d.Count-1)
	return math.Abs(value-d.Mean) > d.Threshold*math.Sqrt(variance)
}
func (d Detector) Variance() float64 {
	if d.Count < 2 {
		return 0
	}
	return d.StdDev / float64(d.Count-1)
}
func (d Detector) ZScore(value float64) float64 {
	variance := d.Variance()
	if variance == 0 {
		return 0
	}
	return (value - d.Mean) / math.Sqrt(variance)
}
