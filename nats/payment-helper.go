package nats

type Status uint

const (
	AwaitingPaymentConfirmation Status = 1
	PaymentConfirmed            Status = 2
	PaymentRejected             Status = 3
)
