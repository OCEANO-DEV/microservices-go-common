package models

import (
	"crypto/ecdsa"
	"time"
)

type PublicKey struct {
	Key       *ecdsa.PublicKey `json:"key"`
	Kid       string           `json:"kid"`
	ExpiresAt time.Time        `json:"expires_at"`
}
