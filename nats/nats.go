package nats

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/services"
)

func NewNats(config *config.Config, service services.CertificatesService) (*nats.Conn, error) {
	caCertBytes, err := service.GetCertificateCA()
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCertBytes)

	tls := &tls.Config{
		MinVersion: tls.VersionTLS12,
		// InsecureSkipVerify: true,
		GetCertificate: service.GetLocalCertificate,
		RootCAs:        pool,
	}
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
		nats.Secure(tls),
	)

	return nc, err
}

func NewJetStream(nc *nats.Conn, streamName string, subjects []string) (nats.JetStreamContext, error) {
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	stream, _ := js.StreamInfo(streamName)
	if stream == nil {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: subjects,
		})
		if err != nil {
			return nil, err
		}
	} else {
		_, err = js.UpdateStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: subjects,
		})
		if err != nil {
			return nil, err
		}
	}

	return js, nil
}
