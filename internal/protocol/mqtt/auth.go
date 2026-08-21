package mqtt

import (
	"errors"
	"fmt"
	"strings"
)

var ErrUnauthorized = errors.New("mqtt publish unauthorized")

type ACL struct{ Allow map[string][]string }

func (a ACL) Allowed(client, topic string) bool {
	for _, pattern := range a.Allow[client] {
		if match(pattern, topic) {
			return true
		}
	}
	return false
}
func (a ACL) Authorize(client, topic string) error {
	if strings.TrimSpace(client) == "" || !a.Allowed(client, topic) {
		message := fmt.Sprintf("mqtt acl rejected %q: %v", client, ErrUnauthorized)
		return fmt.Errorf("%s", message)
	}
	return nil
}
func NormalizeTopic(topic string) string {
	topic = strings.TrimSpace(topic)
	topic = strings.TrimPrefix(topic, "/")
	return topic
}
func IsSystemTopic(topic string) bool { return strings.HasPrefix(topic, "$SYS/") }
