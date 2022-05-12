package security

import (
	"fmt"
	"time"

	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/models"
	"github.com/oceano-dev/microservices-go-common/services"
)

type managerSecurityKeys struct {
	config  *config.Config
	service services.SecurityKeysService
}

var (
	publicKeys        []*models.PublicKey
	refreshPublicKeys = time.Now().UTC()
)

func NewManagerSecurityKeys(
	config *config.Config,
	service services.SecurityKeysService,
) *managerSecurityKeys {
	return &managerSecurityKeys{
		config:  config,
		service: service,
	}
}

func (m *managerSecurityKeys) GetAllPublicKeys() []*models.PublicKey {
	if publicKeys == nil {
		m.refreshPublicKeys()
	}

	publicKeysRefresh := refreshPublicKeys.Sub(time.Now().UTC()) <= 0
	if publicKeysRefresh {
		m.refreshPublicKeys()
		fmt.Println("refresh public keys")
	}

	return publicKeys
}

func (m *managerSecurityKeys) refreshPublicKeys() {
	newestPublicKeys, err := m.service.GetAllPublicKeys()
	if err != nil {
		fmt.Println(err)
	}

	publicKeys = nil
	publicKeys = append(publicKeys, newestPublicKeys...)
	refreshPublicKeys = time.Now().UTC().Add(time.Minute * time.Duration(m.config.SecurityKeys.MinutesToRefreshPublicKeys))
}
