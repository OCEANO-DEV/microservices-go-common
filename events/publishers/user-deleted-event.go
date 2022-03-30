package publishers

import "github.com/oceano-dev/microservices-go-common/events"

type Publishers interface {
	UserDeletedEvent(subject events.Subject, data []byte) error
}

type publishers struct {
	publisher events.Publisher
}

func (p publishers) UserDeletedEvent(subject events.Subject, data []byte) error {
	err := p.publisher.Publish(events.CustomerDeleted, data)
	if err != nil {
		return err
	}

	return nil
}
