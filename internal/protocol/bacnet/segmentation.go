package bacnet

type Segmenter struct{ Max int }

func (s Segmenter) Split(data []byte) [][]byte {
	size := s.Max
	if size <= 0 {
		size = 480
	}
	out := [][]byte{}
	for len(data) > 0 {
		n := size
		if n > len(data) {
			n = len(data)
		}
		out = append(out, append([]byte(nil), data[:n]...))
		data = data[n:]
	}
	return out
}
func (s Segmenter) Join(parts [][]byte) []byte {
	size := 0
	for _, part := range parts {
		size += len(part)
	}
	out := make([]byte, 0, size)
	for _, part := range parts {
		out = append(out, part...)
	}
	return out
}
func (s Segmenter) Count(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	size := s.Max
	if size <= 0 {
		size = 480
	}
	return (len(data) + size - 1) / size
}
