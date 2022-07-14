package models

type Status uint

const (
	SentForPaymentConfirmation  Status = 0
	AwaitingPaymentConfirmation Status = 1
	PaymentConfirmed            Status = 2
	PaymentRejected             Status = 3
	PaymentCanceled             Status = 4
	OrderCreated                Status = 5
	OrderCanceled               Status = 6
)

func (status Status) String() string {
	switch status {
	case SentForPaymentConfirmation:
		return "Sent for payment confirmation"
	case AwaitingPaymentConfirmation:
		return "Awaiting payment confirmation"
	case PaymentConfirmed:
		return "Payment confirmed"
	case PaymentRejected:
		return "Payment rejected"
	case PaymentCanceled:
		return "Payment canceled"
	case OrderCreated:
		return "Order created"
	case OrderCanceled:
		return "Order canceled"
	}

	return "unknown"
}
