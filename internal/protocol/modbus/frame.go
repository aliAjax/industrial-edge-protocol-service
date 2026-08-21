package modbus

import (
	"encoding/binary"
	"errors"
)

var ErrFrame = errors.New("invalid modbus tcp frame")

type Function uint8

const (
	ReadCoils              Function = 1
	ReadDiscreteInputs     Function = 2
	ReadHoldingRegisters   Function = 3
	ReadInputRegisters     Function = 4
	WriteSingleRegister    Function = 6
	WriteMultipleRegisters Function = 16
)

type Frame struct {
	Transaction uint16
	Unit        uint8
	Function    Function
	Payload     []byte
}

func Decode(data []byte) (Frame, error) {
	if len(data) < 8 {
		return Frame{}, ErrFrame
	}
	length := int(binary.BigEndian.Uint16(data[4:6]))
	if length < 2 || length+6 != len(data) || length > 253 {
		return Frame{}, ErrFrame
	}
	payload := make([]byte, len(data)-8)
	copy(payload, data[8:])
	f := Frame{
		Transaction: binary.BigEndian.Uint16(data[:2]),
		Unit:        data[6],
		Function:    Function(data[7]),
		Payload:     payload,
	}
	if f.Function&0x80 != 0 {
		if len(f.Payload) != 1 {
			return Frame{}, ErrFrame
		}
	}
	return f, nil
}
func Encode(f Frame) ([]byte, error) {
	if len(f.Payload) > 245 {
		return nil, ErrFrame
	}
	b := make([]byte, 8+len(f.Payload))
	binary.BigEndian.PutUint16(b[:2], f.Transaction)
	binary.BigEndian.PutUint16(b[2:4], 0)
	binary.BigEndian.PutUint16(b[4:6], uint16(2+len(f.Payload)))
	b[6] = f.Unit
	b[7] = byte(f.Function)
	copy(b[8:], f.Payload)
	return b, nil
}
func ReadRequest(transaction uint16, unit uint8, function Function, address, count uint16) ([]byte, error) {
	if function < ReadCoils || function > ReadInputRegisters || count == 0 || count > 125 {
		return nil, ErrFrame
	}
	p := make([]byte, 4)
	binary.BigEndian.PutUint16(p[:2], address)
	binary.BigEndian.PutUint16(p[2:], count)
	return Encode(Frame{Transaction: transaction, Unit: unit, Function: function, Payload: p})
}
func ParseReadResponse(f Frame) ([]uint16, error) {
	if len(f.Payload) < 1 {
		return nil, ErrFrame
	}
	n := int(f.Payload[0])
	if n%2 != 0 || n+1 != len(f.Payload) {
		return nil, ErrFrame
	}
	out := make([]uint16, n/2)
	for i := range out {
		out[i] = binary.BigEndian.Uint16(f.Payload[1+i*2:])
	}
	return out, nil
}
func Exception(code byte) error { return errors.New("modbus exception code") }
