package amqp

import "github.com/streadway/amqp"

const (
	RetrySuffix      = ".retry"
	QueueRetrySuffix = ".retry"
	delay            = 3
)

type IAMQPConnection interface {
	Channel() (IAMQPChannel, error)
	Close() error
}

type IAMQPChannel interface {
	IsClosed() bool
	Close() error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	DeclareExchange(exchangeName string, exchangeType string) (err error)
	DeclareQueue(queue, queueType, exchangeName, routeKey string, interval int) (amqp.Queue, error)
	BindQueue(q amqp.Queue, routeKey, exchangeName string) error
	NewQueue(exchangeName, queueName, queueType, routeKey string, interval int) (amqp.Queue, error)
	ReadMessage(q amqp.Queue) (<-chan amqp.Delivery, error)
	PublishMessage(q *amqp.Queue, routingKey, contentType, exchange string, body []byte) error
	GetChannel() *amqp.Channel
}
