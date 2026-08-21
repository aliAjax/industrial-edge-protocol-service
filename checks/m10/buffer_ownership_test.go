package m10_test

import (
	"bytes"
	"industrial-edge-protocol/internal/protocol/modbus"
	"industrial-edge-protocol/internal/transport"
	"testing"
)

func TestM10Frame(t *testing.T) {
	w, _ := modbus.Encode(modbus.Frame{Transaction: 7, Unit: 2, Function: modbus.ReadHoldingRegisters, Payload: []byte{2, 0, 9}})
	f, e := modbus.Decode(w)
	if e != nil {
		t.Fatal(e)
	}
	w[8] = 99
	if f.Payload[0] != 2 {
		t.Fatal("frame alias")
	}
}
func TestM10Request(t *testing.T) {
	p := []byte{0, 12, 0, 3}
	f := modbus.Frame{Unit: 1, Function: modbus.ReadHoldingRegisters, Payload: p}
	w := append([]byte(nil), p...)
	r, e := modbus.ParseRequest(f)
	if e != nil || r.Address != 12 || r.Quantity != 3 || !bytes.Equal(f.Payload, w) {
		t.Fatal("request alias")
	}
}
func TestM10Register(t *testing.T) {
	p := []byte{0, 0, 128, 63}
	w := append([]byte(nil), p...)
	_ = modbus.Float32(p, modbus.LittleEndian)
	if !bytes.Equal(p, w) {
		t.Fatal("register alias")
	}
}
func TestM10Chunk(t *testing.T) {
	p := []byte{1, 2, 3, 4}
	c := transport.Chunk(p, 2)
	if len(c) != 2 {
		t.Fatal("count")
	}
	c[0][0] = 9
	if p[0] != 1 {
		t.Fatal("chunk alias")
	}
}
