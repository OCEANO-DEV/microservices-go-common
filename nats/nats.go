package nats

import (
	"log"
	"time"

	"github.com/oceano-dev/microservices-go-common/config"

	"github.com/nats-io/stan.go"
)

func NewNats(config *config.Config) (stan.Conn, error) {
	return stan.Connect(
		config.Nats.ClusterId,
		config.Nats.ClientId,
		stan.ConnectWait(time.Second*time.Duration(config.Nats.ConnectWait)),
		stan.PubAckWait(time.Second*time.Duration(config.Nats.PubAckWait)),
		stan.NatsURL(config.Nats.Url),
		stan.Pings(config.Nats.Interval, config.Nats.MaxOut),
		stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
			log.Fatalf("Connection lost: %v", err)
		}),
	)
}
