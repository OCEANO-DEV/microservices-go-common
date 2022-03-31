package publishers

import (
	"github.com/oceano-dev/microservices-go-common/events"
)

type IUserDeletedEvent interface {
	Publish(data []byte) error
}

type userDeletedEvent struct {
	publisher events.IPublisher
}

func (u *userDeletedEvent) Publish(data []byte) error {
	err := u.publisher.Publish(events.CustomerDeleted, data)
	if err != nil {
		return err
	}

	return nil
}
