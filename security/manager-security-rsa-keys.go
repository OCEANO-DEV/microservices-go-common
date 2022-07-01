package security

import (
	"fmt"
	"time"

	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/models"
	"github.com/oceano-dev/microservices-go-common/services"
)

type managerSecurityRSAKeys struct {
	config  *config.Config
	service services.SecurityRSAKeysService
}

var (
	rsaPublicKeys        []*models.RSAPublicKey
	refreshRSAPublicKeys = time.Now().UTC()
)

func NewManagerSecurityRSAKeys(
	config *config.Config,
	service services.SecurityRSAKeysService,
) *managerSecurityRSAKeys {
	return &managerSecurityRSAKeys{
		config:  config,
		service: service,
	}
}

func (m *managerSecurityRSAKeys) GetAllRSAPublicKeys() []*models.RSAPublicKey {
	if publicKeys == nil {
		m.refreshRSAPublicKeys()
	}

	rsaPublicKeysRefresh := refreshRSAPublicKeys.Before(time.Now().UTC())
	if rsaPublicKeysRefresh {
		m.refreshRSAPublicKeys()
		fmt.Println("refresh RSA public keys")
	}

	return rsaPublicKeys
}

func (m *managerSecurityRSAKeys) refreshRSAPublicKeys() {
	newestRSAPublicKeys, err := m.service.GetAllRSAPublicKeys()
	if err != nil {
		fmt.Println(err)
	}

	rsaPublicKeys = nil
	rsaPublicKeys = append(rsaPublicKeys, newestRSAPublicKeys...)
	refreshRSAPublicKeys = time.Now().UTC().Add(time.Minute * time.Duration(m.config.SecurityRSAKeys.MinutesToRefreshRSAPublicKeys))
}
