package security

import (
	"fmt"
	"time"

	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/models"
	"github.com/oceano-dev/microservices-go-common/services"
)

type ManagerSecurityKeys struct {
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
) *ManagerSecurityKeys {
	return &ManagerSecurityKeys{
		config:  config,
		service: service,
	}
}

func (m *ManagerSecurityKeys) GetAllPublicKeys() []*models.PublicKey {
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

func (m *ManagerSecurityKeys) refreshPublicKeys() {
	newestPublicKeys, err := m.service.GetAllPublicKeys()
	if err != nil {
		fmt.Println(err)
	}

	publicKeys = nil
	publicKeys = append(publicKeys, newestPublicKeys...)
	refreshPublicKeys = time.Now().UTC().Add(time.Minute * time.Duration(m.config.SecurityKeys.MinutesToRefreshPublicKeys))
}
