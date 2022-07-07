package nats

type CustomerSubject string
type OrderSubject string
type ProductSubject string
type StoreSubject string
type UserSubject string

const (
	CustomerDeleted       CustomerSubject = "customer:deleted"
	OrderCreate           OrderSubject    = "order:create"
	OrderPayment          OrderSubject    = "order:payment"
	OrderStatus           OrderSubject    = "order:status"
	ProductCreateMongo    ProductSubject  = "product:create-mongo"
	ProductCreatePostgres ProductSubject  = "product:create-postgres"
	StoreBook             StoreSubject    = "store:storebook"
	StoreBookMongo        StoreSubject    = "store:storebook-mongo"
	StoreBooked           StoreSubject    = "store:storebooked"
	StoreCreateMongo      StoreSubject    = "store:storecreate-mongo"
	StoreCreatePostgres   StoreSubject    = "store:storecreate-postgres"
	StorePaid             StoreSubject    = "store:storepaid"
	StorePayment          StoreSubject    = "store:storepayment"
	StorePaymentMongo     StoreSubject    = "store:storepayment-mongo"
	StorePaymentPostgres  StoreSubject    = "store:storepayment-postgres"
	StoreUnbookMongo      StoreSubject    = "store:storeunbook-mongo"
	StoreUnbookPostgres   StoreSubject    = "store:storeunbook-postgres"
	UserDeleted           UserSubject     = "user:deleted"
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

func GetProductSubjects() []string {
	return []string{
		string(ProductCreateMongo),
		string(ProductCreatePostgres),
	}
}

func GetStoreSubjects() []string {
	return []string{
		string(StoreBook),
		string(StoreBookMongo),
		string(StoreBooked),
		string(StoreCreateMongo),
		string(StoreCreatePostgres),
		string(StorePaid),
		string(StorePayment),
		string(StorePaymentMongo),
		string(StorePaymentPostgres),
		string(StoreUnbookMongo),
		string(StoreUnbookPostgres),
	}
}

func GetUserSubjects() []string {
	return []string{
		string(UserDeleted),
	}
}
