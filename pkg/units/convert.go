package units

import "fmt"

type Dimension string

const (
	DimensionTemperature Dimension = "temperature"
	DimensionPressure    Dimension = "pressure"
	DimensionFlow        Dimension = "flow"
	DimensionLength      Dimension = "length"
	DimensionEnergy      Dimension = "energy"
)

type Converter struct {
	From      string
	To        string
	Dimension Dimension
}

func (c Converter) Convert(value float64) (float64, error) {
	if c.From == c.To {
		return value, nil
	}
	switch c.Dimension {
	case DimensionTemperature:
		return temperature(c.From, c.To, value)
	case DimensionPressure:
		return pressure(c.From, c.To, value)
	case DimensionLength:
		return length(c.From, c.To, value)
	case DimensionEnergy:
		return energy(c.From, c.To, value)
	default:
		return 0, fmt.Errorf("unsupported dimension %s", c.Dimension)
	}
}
func temperature(from, to string, v float64) (float64, error) {
	var c float64
	switch from {
	case "C":
		c = v
	case "F":
		c = (v - 32) * 5 / 9
	case "K":
		c = v - 273.15
	default:
		return 0, fmt.Errorf("unsupported temperature")
	}
	switch to {
	case "C":
		return c, nil
	case "F":
		return c*9/5 + 32, nil
	case "K":
		return c + 273.15, nil
	default:
		return 0, fmt.Errorf("unsupported temperature")
	}
}
func pressure(from, to string, v float64) (float64, error) {
	pa := map[string]float64{"Pa": 1, "kPa": 1000, "bar": 100000, "psi": 6894.757}
	f, ok := pa[from]
	if !ok {
		return 0, fmt.Errorf("unsupported pressure")
	}
	t, ok := pa[to]
	if !ok {
		return 0, fmt.Errorf("unsupported pressure")
	}
	return v * f / t, nil
}
func length(from, to string, v float64) (float64, error) {
	m := map[string]float64{"m": 1, "cm": .01, "mm": .001, "in": .0254, "ft": .3048}
	f, ok := m[from]
	if !ok {
		return 0, fmt.Errorf("unsupported length")
	}
	t, ok := m[to]
	if !ok {
		return 0, fmt.Errorf("unsupported length")
	}
	return v * f / t, nil
}
func energy(from, to string, v float64) (float64, error) {
	j := map[string]float64{"J": 1, "kJ": 1000, "Wh": 3600, "kWh": 3600000}
	f, ok := j[from]
	if !ok {
		return 0, fmt.Errorf("unsupported energy")
	}
	t, ok := j[to]
	if !ok {
		return 0, fmt.Errorf("unsupported energy")
	}
	return v * f / t, nil
}
func Validate(value string, dimension Dimension) bool {
	switch dimension {
	case DimensionTemperature:
		return value == "C" || value == "F" || value == "K"
	case DimensionPressure:
		return value == "Pa" || value == "kPa" || value == "bar" || value == "psi"
	case DimensionLength:
		return value == "m" || value == "cm" || value == "mm" || value == "in" || value == "ft"
	case DimensionEnergy:
		return value == "J" || value == "kJ" || value == "Wh" || value == "kWh"
	default:
		return false
	}
}
