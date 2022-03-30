package events

type Subject string

const (
	CustomerCreated Subject = "customer:created"
	CustomerUpdated Subject = "customer:updated"
	CustomerDeleted Subject = "customer:deleted"
)
