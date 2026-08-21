package modbus

import "fmt"

type Request struct {
	Unit     uint8
	Function Function
	Address  uint16
	Quantity uint16
}

func (r Request) Validate() error {
	if r.Unit == 0 || r.Unit > 247 {
		return fmt.Errorf("invalid unit")
	}
	if !ValidateAddress(r.Address, r.Quantity) {
		return ErrFrame
	}
	switch r.Function {
	case ReadCoils, ReadDiscreteInputs, ReadHoldingRegisters, ReadInputRegisters:
		if r.Quantity > 125 {
			return ErrFrame
		}
	case WriteMultipleRegisters:
		if r.Quantity > 123 {
			return ErrFrame
		}
	default:
		return fmt.Errorf("unsupported function")
	}
	return nil
}
func (r Request) Encode(transaction uint16) ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	p := make([]byte, 4)
	p[0] = byte(r.Address >> 8)
	p[1] = byte(r.Address)
	p[2] = byte(r.Quantity >> 8)
	p[3] = byte(r.Quantity)
	return Encode(Frame{Transaction: transaction, Unit: r.Unit, Function: r.Function, Payload: p})
}
func ParseRequest(f Frame) (Request, error) {
	if len(f.Payload) != 4 {
		return Request{}, ErrFrame
	}
	r := Request{Unit: f.Unit, Function: f.Function, Address: uint16(f.Payload[0])<<8 | uint16(f.Payload[1]), Quantity: uint16(f.Payload[2])<<8 | uint16(f.Payload[3])}
	return r, r.Validate()
}
