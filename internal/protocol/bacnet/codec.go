package bacnet

import (
	"encoding/binary"
	"errors"
	"fmt"
)

var ErrMalformedAPDU = errors.New("malformed bacnet apdu")

type Header struct {
	Version     byte
	MessageType byte
	Length      uint16
	InvokeID    byte
}
type APDU struct {
	Header  Header
	Payload []byte
}

func Decode(data []byte) (APDU, error) {
	if len(data) < 5 {
		return APDU{}, fmt.Errorf("short bacnet packet: %w", ErrMalformedAPDU)
	}
	length := binary.BigEndian.Uint16(data[2:4])
	if int(length) != len(data) {
		return APDU{}, fmt.Errorf("invalid bacnet length: %w", ErrMalformedAPDU)
	}
	return APDU{Header: Header{Version: data[0], MessageType: data[1], Length: length, InvokeID: data[4]}, Payload: append([]byte(nil), data[5:]...)}, nil
}
func Encode(p APDU) []byte {
	p.Header.Length = uint16(5 + len(p.Payload))
	out := make([]byte, p.Header.Length)
	out[0] = p.Header.Version
	out[1] = p.Header.MessageType
	binary.BigEndian.PutUint16(out[2:4], p.Header.Length)
	out[4] = p.Header.InvokeID
	copy(out[5:], p.Payload)
	return out
}

type ObjectID struct {
	Type     uint16
	Instance uint32
}

func (o ObjectID) Encode() []byte {
	value := uint32(o.Type)<<22 | o.Instance&0x3fffff
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, value)
	return b
}
func DecodeObjectID(data []byte) (ObjectID, error) {
	if len(data) < 4 {
		return ObjectID{}, fmt.Errorf("short object id")
	}
	value := binary.BigEndian.Uint32(data)
	return ObjectID{Type: uint16(value >> 22), Instance: value & 0x3fffff}, nil
}

type Property struct {
	Object     ObjectID
	PropertyID uint32
	ArrayIndex *uint32
	Value      []byte
}

func EncodeProperty(p Property) []byte {
	out := append([]byte{}, p.Object.Encode()...)
	id := make([]byte, 4)
	binary.BigEndian.PutUint32(id, p.PropertyID)
	out = append(out, id...)
	if p.ArrayIndex != nil {
		out = append(out, 1)
		binary.BigEndian.PutUint32(id, *p.ArrayIndex)
		out = append(out, id...)
	} else {
		out = append(out, 0)
	}
	out = append(out, p.Value...)
	return out
}
func DecodeProperty(data []byte) (Property, error) {
	if len(data) < 9 {
		return Property{}, fmt.Errorf("short property")
	}
	object, err := DecodeObjectID(data[:4])
	if err != nil {
		return Property{}, err
	}
	p := Property{Object: object, PropertyID: binary.BigEndian.Uint32(data[4:8])}
	offset := 9
	if data[8] == 1 {
		if len(data) < 13 {
			return Property{}, fmt.Errorf("short array index")
		}
		index := binary.BigEndian.Uint32(data[9:13])
		p.ArrayIndex = &index
		offset = 13
	}
	if offset > len(data) {
		return Property{}, fmt.Errorf("invalid property")
	}
	p.Value = append([]byte(nil), data[offset:]...)
	return p, nil
}

type Priority uint8

const (
	Normal Priority = iota
	Urgent
	Critical
)

func ValidatePriority(p Priority) bool     { return p <= Critical }
func IsConfirmedService(service byte) bool { return service == 5 || service == 6 }
func APDUType(value byte) string {
	switch value {
	case 0:
		return "confirmed-request"
	case 1:
		return "unconfirmed-request"
	case 2:
		return "simple-ack"
	case 3:
		return "complex-ack"
	case 4:
		return "segment-ack"
	case 5:
		return "error"
	default:
		return "unknown"
	}
}
