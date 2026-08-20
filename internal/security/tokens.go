package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
	"time"
)

var ErrToken = errors.New("invalid capability token")

type Token struct {
	Subject   string
	Scope     []string
	ExpiresAt time.Time
	Nonce     string
	Digest    string
}

func Issue(subject string, scope []string, ttl time.Duration) (Token, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return Token{}, err
	}
	nonce := base64.RawURLEncoding.EncodeToString(b)
	t := Token{Subject: subject, Scope: append([]string(nil), scope...), ExpiresAt: time.Now().Add(ttl), Nonce: nonce}
	t.Digest = Digest(t)
	return t, nil
}
func Digest(t Token) string {
	sum := sha256.Sum256([]byte(t.Subject + "|" + strings.Join(t.Scope, ",") + "|" + t.Nonce))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
func (t Token) Valid(now time.Time) bool {
	return t.Subject != "" && !t.ExpiresAt.Before(now) && t.Digest == Digest(t)
}
func (t Token) Allows(scope string) bool {
	for _, v := range t.Scope {
		if v == scope || v == "*" {
			return true
		}
	}
	return false
}
