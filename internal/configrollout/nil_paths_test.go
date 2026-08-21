package configrollout_test

import (
	"context"
	"errors"
	"testing"

	"industrial-edge-protocol/internal/configrollout"
	"industrial-edge-protocol/internal/domain"
	"industrial-edge-protocol/internal/gateway"
)

func TestRolloutZeroCapabilitiesAreInitialized(t *testing.T) {
	compiled, err := configrollout.Compile(3, domain.Gateway{ID: "gw-1"}, nil, configrollout.CapabilitySet{Protocols: []string{"modbus"}})
	if err != nil {
		t.Fatal(err)
	}
	if compiled.Gateway != "gw-1" || compiled.Digest == "" {
		t.Fatalf("invalid compiled config: %+v", compiled)
	}
}

func TestRolloutAckOnZeroValueManager(t *testing.T) {
	var manager configrollout.Manager
	err := manager.Start(domain.Rollout{ID: "rollout-1", Percent: 100, GatewayIDs: []domain.ID{"gw-1"}})
	if err != nil {
		t.Fatal(err)
	}
	manager.Ack("rollout-1", "gw-1", true)
	rollout, ok := manager.Get("rollout-1")
	if !ok || rollout.State != "completed" {
		t.Fatalf("zero-value manager did not complete rollout: %+v", rollout)
	}
}

func TestRolloutEmptySelectionStaysValid(t *testing.T) {
	selected := (configrollout.Strategy{Name: "empty", Percentage: 0}).Select(nil)
	if selected == nil || len(selected) != 0 {
		t.Fatalf("empty rollout selection must be a usable empty slice: %#v", selected)
	}
}

func TestRolloutRejectsTypedNilTransport(t *testing.T) {
	var memory *gateway.MemoryTransport
	var transport gateway.Transport = memory
	err := transport.ApplyConfig(context.Background(), domain.Gateway{ID: "gw-1"}, []byte("{}"))
	if !errors.Is(err, gateway.ErrTransportClosed) {
		t.Fatalf("typed nil transport was not rejected: %v", err)
	}
}
