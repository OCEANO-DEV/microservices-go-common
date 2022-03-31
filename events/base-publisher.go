package events

import (
	"github.com/nats-io/stan.go"
)

type IPublisher interface {
	Publish(subject Subject, data []byte) error
}

type publisher struct {
	stan stan.Conn
}

func NewPublisher(
	stan stan.Conn,
) *publisher {
	return &publisher{
		stan: stan,
	}
}

func (p *publisher) Publish(subject Subject, data []byte) error {
	err := p.stan.Publish(string(subject), data)
	if err != nil {
		return err
	}

	return nil
}
