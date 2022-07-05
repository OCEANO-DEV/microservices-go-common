package nats

type Subject string

const (
	CustomerDeleted Subject = "customer:deleted"
	OrderCreate     Subject = "order:create"
	OrderPayment    Subject = "order:payment"
	OrderStatus     Subject = "order:status"
	StoreBook       Subject = "store:book"
	StoreBooked     Subject = "store:booked"
	StorePaid       Subject = "store:paid"
	StorePayment    Subject = "store:payment"
	UserDeleted     Subject = "user:deleted"
)

func GetSubjects() []string {
	return []string{
		string(CustomerDeleted),
		string(OrderCreate),
		string(OrderPayment),
		string(OrderStatus),
		string(StoreBook),
		string(StoreBooked),
		string(StorePaid),
		string(StorePayment),
		string(UserDeleted),
	}
}
