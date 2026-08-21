package modbus

import (
	"bytes"
	"testing"

	"industrial-edge-protocol/internal/transport"
)

func TestModbusDecodedFrameOwnsPayload(t *testing.T) {
	wire, err := Encode(Frame{Transaction: 7, Unit: 2, Function: ReadHoldingRegisters, Payload: []byte{2, 0, 9}})
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := Decode(wire)
	if err != nil {
		t.Fatal(err)
	}
	wire[8] = 99
	if decoded.Payload[0] != 2 {
		t.Fatalf("decoded payload changed after receive buffer reuse: %v", decoded.Payload)
	}
}

func TestModbusRequestSurvivesBufferReuse(t *testing.T) {
	payload := []byte{0, 12, 0, 3}
	frame := Frame{Unit: 1, Function: ReadHoldingRegisters, Payload: payload}
	before := append([]byte(nil), payload...)
	request, err := ParseRequest(frame)
	if err != nil {
		t.Fatal(err)
	}
	if request.Address != 12 || request.Quantity != 3 {
		t.Fatalf("unexpected request: %+v", request)
	}
	if !bytes.Equal(frame.Payload, before) {
		t.Fatalf("request parsing mutated the frame payload: %v", frame.Payload)
	}
}

func TestModbusRegisterDecodeDoesNotMutateInput(t *testing.T) {
	input := []byte{0, 0, 128, 63}
	want := append([]byte(nil), input...)
	_ = Float32(input, LittleEndian)
	if !bytes.Equal(input, want) {
		t.Fatalf("register decoder rewrote the caller buffer: %v", input)
	}
}

func TestModbusChunkBoundaryUsesIndependentSlices(t *testing.T) {
	input := []byte{1, 2, 3, 4}
	chunks := transport.Chunk(input, 2)
	if len(chunks) != 2 {
		t.Fatalf("unexpected chunk count: %d", len(chunks))
	}
	chunks[0][0] = 9
	if input[0] != 1 {
		t.Fatalf("chunk mutation leaked into input: %v", input)
	}
}
