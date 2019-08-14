package token

import (
	"time"
)

// Provider is for providing new token
type Provider interface {
	Store(data map[string]interface{}, exp time.Duration) (string, error)

	Fetch(token string) (map[string]interface{}, error)

	Delete(token string)
}
