package events

import "github.com/nats-io/stan.go"

type Publisher interface {
	Publish(subject Subject, data []byte) error
}

type publisher struct {
	client stan.Conn
}

func NewPublisher(
	client stan.Conn,
) *publisher {
	return &publisher{
		client: client,
	}
}

func (p *publisher) Publish(subject Subject, data []byte) error {
	err := p.client.Publish(string(subject), data)
	if err != nil {
		return err
	}

	return nil
}
