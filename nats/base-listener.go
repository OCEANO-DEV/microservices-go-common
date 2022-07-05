package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	ackWait = 5 * time.Second
)

type Listener interface {
	Listener(subject string, queueGroupName string, handler nats.MsgHandler)
}

type listener struct {
	js nats.JetStream
}

func NewListener(
	js nats.JetStream,
) *listener {
	return &listener{
		js: js,
	}
}

func (l *listener) Listener(subject string, queueGroupName string, handler nats.MsgHandler) {
	_, err := l.js.QueueSubscribe(
		subject,
		queueGroupName,
		handler,
		nats.Durable(queueGroupName),
		nats.DeliverAll(),
		nats.ManualAck(),
		nats.AckWait(ackWait),
	)
	if err != nil {
		fmt.Println(fmt.Errorf("Subject: %v, QueueSubscribe: %v, Error: %v", subject, queueGroupName, err))
	}
}
