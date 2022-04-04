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
	stan stan.Conn
}

func NewListener(
	stan stan.Conn,
) *listener {
	return &listener{
		stan: stan,
	}
}

func (l *listener) Listener(subject Subject, queueGroupName string, handler stan.MsgHandler) {
	_, err := l.stan.QueueSubscribe(
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
		if err := l.stan.Close(); err != nil {
			fmt.Println(fmt.Errorf("Subject: %v, conn.Close error: %v", string(subject), err))
		}
	}
}
