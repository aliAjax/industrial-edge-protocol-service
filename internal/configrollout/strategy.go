package configrollout

import (
	"crypto/sha256"
	"encoding/hex"
	"industrial-edge-protocol/internal/domain"
	"sort"
)

type Strategy struct {
	Name       string
	Percentage int
	GatewayIDs []domain.ID
	Canary     bool
}

func (s Strategy) Validate() bool { return s.Name != "" && s.Percentage >= 0 && s.Percentage <= 100 }
func (s Strategy) Select(all []domain.Gateway) []domain.Gateway {
	var selected []domain.Gateway
	if len(s.GatewayIDs) > 0 {
		wanted := map[domain.ID]bool{}
		for _, id := range s.GatewayIDs {
			wanted[id] = true
		}
		for _, g := range all {
			if wanted[g.ID] {
				selected = append(selected, g)
			}
		}
		return selected
	}
	sort.SliceStable(all, func(i, j int) bool { return all[i].ID < all[j].ID })
	count := len(all) * s.Percentage / 100
	if s.Canary && count < 1 && len(all) > 0 {
		count = 1
	}
	if count > len(all) {
		count = len(all)
	}
	return append(selected, all[:count]...)
}
func StableHash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
func RolloutDigest(version int64, gatewayIDs []domain.ID) string {
	ids := append([]domain.ID(nil), gatewayIDs...)
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	value := StableHash(string(rune(version)))
	for _, id := range ids {
		value = StableHash(value + string(id))
	}
	return value
}
func Difference(old, new []domain.ID) []domain.ID {
	present := map[domain.ID]bool{}
	for _, id := range old {
		present[id] = true
	}
	out := []domain.ID{}
	for _, id := range new {
		if !present[id] {
			out = append(out, id)
		}
	}
	return out
}
