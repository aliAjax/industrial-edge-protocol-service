package mqtt

import (
	"encoding/binary"
	"errors"
)

var ErrPacket = errors.New("invalid mqtt packet")

type Packet struct {
	Type    byte
	Flags   byte
	Payload []byte
}

func Decode(data []byte) (Packet, error) {
	if len(data) < 2 {
		return Packet{}, ErrPacket
	}
	length, n := readLength(data[1:])
	if n == 0 || 2+n+length != len(data) {
		return Packet{}, ErrPacket
	}
	return Packet{Type: data[0] >> 4, Flags: data[0] & 15, Payload: append([]byte(nil), data[1+n:]...)}, nil
}
func Encode(p Packet) []byte {
	length := len(p.Payload)
	head := []byte{p.Type<<4 | p.Flags}
	for {
		digit := byte(length % 128)
		length /= 128
		if length > 0 {
			digit |= 128
		}
		head = append(head, digit)
		if length == 0 {
			break
		}
	}
	return append(head, p.Payload...)
}
func readLength(data []byte) (int, int) {
	multiplier := 1
	value := 0
	for i, b := range data {
		value += int(b&127) * multiplier
		if b&128 == 0 {
			return value, i + 1
		}
		multiplier *= 128
		if multiplier > 128*128*128 {
			return 0, 0
		}
	}
	return 0, 0
}
func EncodeUint16(v uint16) []byte {
	out := make([]byte, 2)
	binary.BigEndian.PutUint16(out, v)
	return out
}
