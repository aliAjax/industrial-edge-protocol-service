package mqtt

import (
	"errors"
	"testing"
)

func TestMQTTDecodeErrorKeepsSentinel(t *testing.T) {
	for _, packet := range [][]byte{{0x30}, {0x30, 0x02, 0x01}} {
		_, err := Decode(packet)
		if !errors.Is(err, ErrPacket) {
			t.Fatalf("packet sentinel lost for %x: %v", packet, err)
		}
	}
}

func TestMQTTACLRejectionIsClassified(t *testing.T) {
	acl := ACL{Allow: map[string][]string{"edge-a": {"plant/+/data"}}}
	err := acl.Authorize("edge-b", "plant/boiler/data")
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("authorization sentinel lost: %v", err)
	}
}

func TestMQTTQoSErrorIsNotRetried(t *testing.T) {
	err := ValidateQoS(QoS(7))
	if !errors.Is(err, ErrInvalidQoS) {
		t.Fatalf("qos sentinel lost: %v", err)
	}
}

func TestMQTTRouterPreservesHandlerErrorChain(t *testing.T) {
	sentinel := errors.New("device rejected publish")
	router := NewRouter()
	router.Subscribe("plant/#", func(Message) error { return sentinel })
	errs := router.Publish(Message{Topic: "plant/line-1/data"})
	if len(errs) != 1 || !errors.Is(errs[0], sentinel) {
		t.Fatalf("handler error identity lost: %v", errs)
	}
}
