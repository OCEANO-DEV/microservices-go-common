package models

type OrderStatus uint

const (
	SentForPaymentConfirmation  OrderStatus = 0
	AwaitingPaymentConfirmation OrderStatus = 1
	PaymentConfirmed            OrderStatus = 2
	PaymentRejected             OrderStatus = 3
	OrderCreated                OrderStatus = 4
	OrderCanceled               OrderStatus = 5
)

func (status OrderStatus) String() string {
	switch status {
	case SentForPaymentConfirmation:
		return "Sent for payment confirmation"
	case AwaitingPaymentConfirmation:
		return "Awaiting payment confirmation"
	case PaymentConfirmed:
		return "Payment confirmed"
	case PaymentRejected:
		return "Payment rejected"
	case OrderCreated:
		return "Order created"
	case OrderCanceled:
		return "Order canceled"
	}

	return "unknown"
}
