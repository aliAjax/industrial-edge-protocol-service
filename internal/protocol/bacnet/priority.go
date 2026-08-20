package bacnet

type PriorityArray struct{ values [16]any }

func (p *PriorityArray) Set(priority int, value any) bool {
	if priority < 1 || priority > 16 {
		return false
	}
	p.values[priority-1] = value
	return true
}
func (p *PriorityArray) Clear(priority int) bool {
	if priority < 1 || priority > 16 {
		return false
	}
	p.values[priority-1] = nil
	return true
}
func (p PriorityArray) Effective() (any, bool) {
	for _, value := range p.values {
		if value != nil {
			return value, true
		}
	}
	return nil, false
}
func (p PriorityArray) At(priority int) (any, bool) {
	if priority < 1 || priority > 16 {
		return nil, false
	}
	value := p.values[priority-1]
	return value, value != nil
}
