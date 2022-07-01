package security

import "github.com/oceano-dev/microservices-go-common/models"

type ManagerSecurityRSAKeys interface {
	GetAllRSAPublicKeys() []*models.RSAPublicKey
}
