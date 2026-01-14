package constants

const (
	// User Exchange
	UserExchangeName = "user_events"

	// User Queue Names
	UserCreatedQueueName = "user_created_queue"
	UserUpdatedQueueName = "user_updated_queue"

	// User Routing Keys
	UserCreatedRouteKey = "user.created"
	UserUpdatedRouteKey = "user.updated"

	// Queue Type
	QueueType = "quorum"

	// Queue Interval (in milliseconds)
	UserCreatedQueueInterval = 3000
	UserUpdatedQueueInterval = 3000
)
