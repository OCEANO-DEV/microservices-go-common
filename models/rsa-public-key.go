package models

import (
	"crypto/rsa"
	"time"
)

type RSAPublicKey struct {
	Key       *rsa.PublicKey
	ExpiresAt time.Time `json:"expires_at"`
}
