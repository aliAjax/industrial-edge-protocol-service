package b3_test

import (
	"errors"
	"industrial-edge-protocol/internal/protocol/bacnet"
	"testing"
)

func TestB3Malformed(t *testing.T) {
	_, e := bacnet.Decode([]byte{1, 2, 3})
	if !errors.Is(e, bacnet.ErrMalformedAPDU) {
		t.Fatal(e)
	}
}
func TestB3TypedService(t *testing.T) {
	e := bacnet.NewServiceError("property", "write-access-denied")
	var s bacnet.ServiceError
	if !errors.As(e, &s) {
		t.Fatal(e)
	}
}
func TestB3Priority(t *testing.T) {
	var p bacnet.PriorityArray
	if !errors.Is(p.SetChecked(17, 42), bacnet.ErrInvalidPriority) {
		t.Fatal("sentinel")
	}
}
func TestB3Segment(t *testing.T) {
	_, e := (bacnet.Segmenter{Max: 2}).JoinAPDU([][]byte{{1, 2}, {3}})
	if !errors.Is(e, bacnet.ErrMalformedAPDU) {
		t.Fatal(e)
	}
}
