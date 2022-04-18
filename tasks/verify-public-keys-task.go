package tasks

import (
	"fmt"
	"sync"
	"time"

	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/security"
)

type VerifyPublicKeysTask struct {
	config  *config.Config
	manager *security.ManagerSecurityKeys
}

func NewVerifyPublicKeysTask(
	config *config.Config,
	manager *security.ManagerSecurityKeys,
) *VerifyPublicKeysTask {
	return &VerifyPublicKeysTask{
		config:  config,
		manager: manager,
	}
}

var mux sync.Mutex

func (task *VerifyPublicKeysTask) ReloadPublicKeys() {
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				_ = task.manager.GetAllPublicKeys()

				fmt.Printf("publickeys success refreshed %s\n", time.Now().UTC())
				ticker.Reset(1 * time.Hour)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
