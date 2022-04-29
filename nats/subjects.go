package nats

type Subject string

const (
	CustomerCreated Subject = "customer:created"
	CustomerUpdated Subject = "customer:updated"
	CustomerDeleted Subject = "customer:deleted"
	UserCreated     Subject = "user:created"
	UserUpdated     Subject = "user:updated"
	UserDeleted     Subject = "user:deleted"
)
