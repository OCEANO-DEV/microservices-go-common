package nats

import (
	"github.com/nats-io/stan.go"
)

type Publisher interface {
	Publish(subject string, data []byte) error
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

func (p *publisher) Publish(subject string, data []byte) error {
	err := p.stan.Publish(subject, data)
	if err != nil {
		return err
	}

	return nil
}
