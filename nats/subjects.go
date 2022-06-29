package nats

type Subject string

const (
	CustomerDeleted Subject = "customer:deleted"
	OrderCreate     Subject = "order:create"
	OrderPaid       Subject = "order:paid"
	OrderPayment    Subject = "order:payment"
	StoreBook       Subject = "store:book"
	StoreBooked     Subject = "store:booked"
	StorePaid       Subject = "store:paid"
	StorePayment    Subject = "store:payment"
	UserDeleted     Subject = "user:deleted"
)
