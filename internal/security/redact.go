package security

import "strings"

var sensitive = []string{"password", "secret", "token", "private_key", "certificate"}

func Redact(values map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range values {
		lower := strings.ToLower(k)
		hidden := false
		for _, word := range sensitive {
			if strings.Contains(lower, word) {
				hidden = true
				break
			}
		}
		if hidden {
			out[k] = "[REDACTED]"
		} else {
			out[k] = v
		}
	}
	return out
}
func IsSensitive(name string) bool {
	lower := strings.ToLower(name)
	for _, word := range sensitive {
		if strings.Contains(lower, word) {
			return true
		}
	}
	return false
}
