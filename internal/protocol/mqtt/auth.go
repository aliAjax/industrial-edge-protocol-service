package mqtt

import "strings"

type ACL struct{ Allow map[string][]string }

func (a ACL) Allowed(client, topic string) bool {
	for _, pattern := range a.Allow[client] {
		if match(pattern, topic) {
			return true
		}
	}
	return false
}
func NormalizeTopic(topic string) string {
	topic = strings.TrimSpace(topic)
	topic = strings.TrimPrefix(topic, "/")
	return topic
}
func IsSystemTopic(topic string) bool { return strings.HasPrefix(topic, "$SYS/") }
