package models

import (
	"crypto/ecdsa"
	"time"
)

type PublicKey struct {
	Key       *ecdsa.PublicKey
	Kid       string    `json:"kid"`
	ExpiresAt time.Time `json:"expires_at"`
}
