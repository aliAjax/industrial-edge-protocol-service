package modbus

import (
	"encoding/binary"
	"math"
)

type ByteOrder int

const (
	BigEndian ByteOrder = iota
	LittleEndian
	WordSwap
)

func Uint16(data []byte, order ByteOrder) uint16 {
	if len(data) < 2 {
		return 0
	}
	if order == LittleEndian {
		return binary.LittleEndian.Uint16(data)
	}
	return binary.BigEndian.Uint16(data)
}
func Int16(data []byte, order ByteOrder) int16 { return int16(Uint16(data, order)) }
func Float32(data []byte, order ByteOrder) float32 {
	if len(data) < 4 {
		return 0
	}
	b := append([]byte(nil), data[:4]...)
	if order == LittleEndian {
		for i, j := 0, 3; i < j; i, j = i+1, j-1 {
			b[i], b[j] = b[j], b[i]
		}
	}
	if order == WordSwap {
		b[0], b[2] = b[2], b[0]
		b[1], b[3] = b[3], b[1]
	}
	return math.Float32frombits(binary.BigEndian.Uint32(b))
}
func PackFloat32(v float32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, math.Float32bits(v))
	return b
}
func ValidateAddress(address, count uint16) bool {
	return uint32(address)+uint32(count) <= 65536 && count > 0
}
