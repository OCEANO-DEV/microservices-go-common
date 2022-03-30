package events

import (
	"fmt"
	"time"

	"github.com/nats-io/stan.go"
)

const (
	ackWait = 5 * time.Second
)

type Listener interface {
	Listener(subject Subject, queueGroupName string, handler stan.MsgHandler)
}

type listener struct {
	client stan.Conn
}

func NewListener(
	client stan.Conn,
) *listener {
	return &listener{
		client: client,
	}
}

func (l *listener) Listener(subject Subject, queueGroupName string, handler stan.MsgHandler) {
	_, err := l.client.QueueSubscribe(
		string(subject),
		queueGroupName,
		handler,
		stan.DurableName(queueGroupName),
		stan.DeliverAllAvailable(),
		stan.SetManualAckMode(),
		stan.AckWait(ackWait),
	)
	if err != nil {
		fmt.Println(fmt.Errorf("Subject: %v, QueueSubscribe: %v, Error: %v", string(subject), queueGroupName, err))
		if err := l.client.Close(); err != nil {
			fmt.Println(fmt.Errorf("Subject: %v, conn.Close error: %v", string(subject), err))
		}
	}
}
