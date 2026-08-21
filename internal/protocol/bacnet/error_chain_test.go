package bacnet

import (
	"errors"
	"testing"
)

func TestBACnetDecodePreservesServiceError(t *testing.T) {
	_, err := Decode([]byte{1, 2, 3})
	if !errors.Is(err, ErrMalformedAPDU) {
		t.Fatalf("apdu sentinel lost: %v", err)
	}
}

func TestBACnetReadonlyPropertyClassification(t *testing.T) {
	err := NewServiceError("property", "write-access-denied")
	var serviceErr ServiceError
	if !errors.As(err, &serviceErr) {
		t.Fatalf("typed service error lost: %v", err)
	}
	if serviceErr.Class != "property" || serviceErr.Code != "write-access-denied" {
		t.Fatalf("unexpected service error: %+v", serviceErr)
	}
}

func TestBACnetPriorityFailureWrapsSentinel(t *testing.T) {
	var priorities PriorityArray
	err := priorities.SetChecked(17, 42)
	if !errors.Is(err, ErrInvalidPriority) {
		t.Fatalf("priority sentinel lost: %v", err)
	}
}

func TestBACnetSegmentJoinPropagatesDecodeError(t *testing.T) {
	_, err := (Segmenter{Max: 2}).JoinAPDU([][]byte{{1, 2}, {3}})
	if !errors.Is(err, ErrMalformedAPDU) {
		t.Fatalf("segmented decode error lost: %v", err)
	}
}
