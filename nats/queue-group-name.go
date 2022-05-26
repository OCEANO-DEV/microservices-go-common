package nats

type QueueGroupName string

const (
	AuthenticationsServiceQueueGroupName QueueGroupName = "authentications-service"
	CustomersServiceQueueGroupName       QueueGroupName = "customers-service"
	ProductsServiceQueueGroupName        QueueGroupName = "products-service"
)
