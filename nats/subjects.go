package nats

type Subject string

const (
	CustomerDeleted Subject = "customer:deleted"
	UserDeleted     Subject = "user:deleted"
	ProductPayment  Subject = "product:payment"
	ProductPaid     Subject = "product:paid"
	ProductBought   Subject = "product:bought"
)
