package transport

import (
	"encoding/binary"
	"errors"
)

var ErrPayload = errors.New("invalid transport payload")

func Frame(payload []byte) []byte {
	out := make([]byte, 4+len(payload))
	binary.BigEndian.PutUint32(out, uint32(len(payload)))
	copy(out[4:], payload)
	return out
}
func Unframe(data []byte) ([]byte, error) {
	if len(data) < 4 {
		return nil, ErrPayload
	}
	n := int(binary.BigEndian.Uint32(data))
	if n < 0 || n != len(data)-4 {
		return nil, ErrPayload
	}
	return append([]byte(nil), data[4:]...), nil
}
func Chunk(data []byte, size int) [][]byte {
	if size <= 0 {
		size = 1024
	}
	out := [][]byte{}
	for len(data) > 0 {
		n := size
		if n > len(data) {
			n = len(data)
		}
		out = append(out, data[:n])
		data = data[n:]
	}
	return out
}
