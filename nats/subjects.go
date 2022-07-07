package nats

type CustomerSubject string
type OrderSubject string
type UserSubject string

const (
	CustomerDeleted CustomerSubject = "customer:deleted"
	OrderCreate     OrderSubject    = "order:create"
	OrderPayment    OrderSubject    = "order:payment"
	OrderStatus     OrderSubject    = "order:status"
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

func GetUserSubjects() []string {
	return []string{
		string(UserDeleted),
	}
}
