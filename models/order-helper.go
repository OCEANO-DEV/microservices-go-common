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
