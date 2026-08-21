package configrollout

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"industrial-edge-protocol/internal/domain"
	"sort"
)

type CapabilitySet struct {
	Protocols []string
	MaxPoints int
	Features  []string
}
type Compiled struct {
	Version  int64
	Digest   string
	Gateway  domain.ID
	Points   []domain.Point
	Warnings []string
}

func Compile(version int64, gateway domain.Gateway, points []domain.Point, caps CapabilitySet) (Compiled, error) {
	sort.Slice(points, func(i, j int) bool { return points[i].ID < points[j].ID })
	if len(points) > caps.MaxPoints && caps.MaxPoints > 0 {
		return Compiled{}, fmt.Errorf("point capacity exceeded")
	}
	for _, p := range points {
		if p.DataType == "" {
			return Compiled{}, fmt.Errorf("point %s missing data type", p.ID)
		}
	}
	raw, _ := json.Marshal(struct {
		V int64
		G domain.ID
		P []domain.Point
	}{version, gateway.ID, points})
	sum := sha256.Sum256(raw)
	return Compiled{Version: version, Digest: fmt.Sprintf("%x", sum), Gateway: gateway.ID, Points: points}, nil
}
func Diff(old, new Compiled) []string {
	changes := []string{}
	if old.Digest != new.Digest {
		changes = append(changes, "digest changed")
	}
	if len(old.Points) != len(new.Points) {
		changes = append(changes, "point count changed")
	}
	return changes
}
