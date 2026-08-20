package serialization

import (
	"encoding/json"
	"sort"
)

func CanonicalMap(value map[string]any) []byte {
	keys := make([]string, 0, len(value))
	for k := range value {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	ordered := make(map[string]any, len(value))
	for _, k := range keys {
		ordered[k] = value[k]
	}
	raw, _ := json.Marshal(ordered)
	return raw
}
func CloneMap(value map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range value {
		out[k] = v
	}
	return out
}
func Merge(base, overlay map[string]any) map[string]any {
	out := CloneMap(base)
	for k, v := range overlay {
		if nested, ok := v.(map[string]any); ok {
			if old, ok := out[k].(map[string]any); ok {
				out[k] = Merge(old, nested)
				continue
			}
		}
		out[k] = v
	}
	return out
}
