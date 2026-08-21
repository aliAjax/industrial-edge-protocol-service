package m9_test

import (
	"errors"
	"industrial-edge-protocol/internal/protocol/mqtt"
	"testing"
)

func TestM9PacketIdentity(t *testing.T) {
	for _, p := range [][]byte{{0x30}, {0x30, 0x02, 0x01}} {
		_, e := mqtt.Decode(p)
		if !errors.Is(e, mqtt.ErrPacket) {
			t.Fatal(e)
		}
	}
}
func TestM9ACL(t *testing.T) {
	a := mqtt.ACL{Allow: map[string][]string{"edge-a": {"plant/+/data"}}}
	if !errors.Is(a.Authorize("edge-b", "plant/boiler/data"), mqtt.ErrUnauthorized) {
		t.Fatal("acl")
	}
}
func TestM9QoS(t *testing.T) {
	if !errors.Is(mqtt.ValidateQoS(mqtt.QoS(7)), mqtt.ErrInvalidQoS) {
		t.Fatal("qos")
	}
}
func TestM9Router(t *testing.T) {
	s := errors.New("rejected")
	r := mqtt.NewRouter()
	r.Subscribe("plant/#", func(mqtt.Message) error { return s })
	e := r.Publish(mqtt.Message{Topic: "plant/line-1/data"})
	if len(e) != 1 || !errors.Is(e[0], s) {
		t.Fatal("route")
	}
}
