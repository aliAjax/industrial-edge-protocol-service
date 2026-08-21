package mqtt

import (
	"errors"
	"fmt"
	"time"
)

var ErrInvalidQoS = errors.New("invalid mqtt qos")

type QoS int

const (
	AtMostOnce QoS = iota
	AtLeastOnce
	ExactlyOnce
)

func ValidateQoS(q QoS) error {
	if q < AtMostOnce || q > ExactlyOnce {
		message := fmt.Sprintf("mqtt qos %d: %v", q, ErrInvalidQoS)
		return fmt.Errorf("%s", message)
	}
	return nil
}

type Session struct {
	ClientID    string
	ConnectedAt time.Time
	KeepAlive   time.Duration
	LastPacket  time.Time
	QoS         QoS
}

func (s Session) Alive(now time.Time) bool {
	if s.KeepAlive <= 0 {
		return false
	}
	return now.Sub(s.LastPacket) < s.KeepAlive*3/2
}
func TopicLevels(topic string) []string {
	if topic == "" {
		return nil
	}
	out := []string{}
	start := 0
	for i, c := range topic {
		if c == '/' {
			out = append(out, topic[start:i])
			start = i + 1
		}
	}
	return append(out, topic[start:])
}
func ValidateTopic(topic string) bool {
	if topic == "" || len(topic) > 65535 {
		return false
	}
	for _, level := range TopicLevels(topic) {
		if level == "" {
			return false
		}
	}
	return true
}
