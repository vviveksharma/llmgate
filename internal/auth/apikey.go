package auth

import (
	"crypto/subtle"
	"strings"
)

type APIKeyValidator struct {
	keys map[string]string // key -> tenant ID
}

func NewAPIKeyValidator(keys map[string]string) *APIKeyValidator {
	return &APIKeyValidator{keys: keys}
}

func (v *APIKeyValidator) Validate(key string) (tenantID string, ok bool) {
	key = strings.TrimSpace(key)
	for k, tid := range v.keys {
		if subtle.ConstantTimeCompare([]byte(k), []byte(key)) == 1 {
			return tid, true
		}
	}
	return "", false
}
