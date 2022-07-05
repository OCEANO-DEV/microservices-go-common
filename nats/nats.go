package nats

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/oceano-dev/microservices-go-common/config"
)

func NewNats(config *config.Config) (*nats.Conn, error) {
	nc, err := nats.Connect(
		config.Nats.Url,
		nats.Timeout(time.Second*time.Duration(config.Nats.ConnectWait)),
		nats.PingInterval(time.Second*time.Duration(config.Nats.Interval)),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(5),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Fatalf("Connection lost: %v", err)
		}),
	)

	return nc, err
}

func NewJetStream(nc *nats.Conn, name string, subjects []string) (nats.JetStreamContext, error) {
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     name,
		Subjects: subjects,
	})
	if err != nil {
		return nil, err
	}

	return js, nil
}
