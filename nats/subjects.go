package nats

type CustomerSubject string
type OrderSubject string
type StoreSubject string
type UserSubject string

const (
	CustomerDeleted CustomerSubject = "customer:deleted"
	OrderCreate     OrderSubject    = "order:create"
	OrderPayment    OrderSubject    = "order:payment"
	OrderStatus     OrderSubject    = "order:status"
	StoreBook       StoreSubject    = "store:book"
	StoreBooked     StoreSubject    = "store:booked"
	StorePaid       StoreSubject    = "store:paid"
	StorePayment    StoreSubject    = "store:payment"
	UserDeleted     UserSubject     = "user:deleted"
)

func GetCustomerSubjects() []string {
	return []string{
		string(CustomerDeleted),
	}
}

func GetOrderSubjects() []string {
	return []string{
		string(OrderCreate),
		string(OrderPayment),
		string(OrderStatus),
	}
}

func GetStoreSubjects() []string {
	return []string{
		string(StoreBook),
		string(StoreBooked),
		string(StorePaid),
		string(StorePayment),
	}
}

func GetUserSubjects() []string {
	return []string{
		string(UserDeleted),
	}
}
