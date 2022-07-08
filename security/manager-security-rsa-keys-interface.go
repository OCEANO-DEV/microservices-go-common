package security

import (
	"crypto/rsa"

	"github.com/oceano-dev/microservices-go-common/models"
)

type ManagerSecurityRSAKeys interface {
	GetAllRSAPublicKeys() []*models.RSAPublicKey
	Encrypt(msg string, publicKey *rsa.PublicKey) (string, error)
	Dencrypt(encryptedBytes []byte, privateKey *rsa.PrivateKey) (string, error)
}
