package amqp

import (
	"boilerblade/helper"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/streadway/amqp"
)

type amqpChannel struct {
	*amqp.Channel
	closed int32
}

// IsClosed indicate closed by developer
func (ch *amqpChannel) IsClosed() bool {
	return (atomic.LoadInt32(&ch.closed) == 1)
}

// Close ensure closed flag set
func (ch *amqpChannel) Close() error {
	if ch.IsClosed() {
		return amqp.ErrClosed
	}

	atomic.StoreInt32(&ch.closed, 1)

	return ch.Channel.Close()
}

// Consume warp amqp.Channel.Consume, the returned delivery will end only when channel closed by developer
func (ch *amqpChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	deliveries := make(chan amqp.Delivery)

	go func() {
		for {
			d, err := ch.Channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
			if err != nil {
				helper.LogError("AMQP consume failed", err, queue, map[string]interface{}{
					"queue":     queue,
					"consumer":  consumer,
					"auto_ack":  autoAck,
					"exclusive": exclusive,
				})
				time.Sleep(delay * time.Second)
				continue
			}

			for msg := range d {
				deliveries <- msg
			}

			// sleep before IsClose call. closed flag may not set before sleep.
			time.Sleep(delay * time.Second)

			if ch.IsClosed() {
				break
			}
		}
	}()

	return deliveries, nil
}

func (ch *amqpChannel) DeclareExchange(exchangeName string, exchangeType string) (err error) {
	err = ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		helper.LogError("AMQP exchange declare failed", err, exchangeName, map[string]interface{}{
			"exchange_name": exchangeName,
			"exchange_type": exchangeType,
		})
		return
	}
	err = ch.ExchangeDeclare(
		exchangeName+RetrySuffix,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		helper.LogError("AMQP retry exchange declare failed", err, exchangeName+RetrySuffix, map[string]interface{}{
			"exchange_name": exchangeName + RetrySuffix,
			"exchange_type": exchangeType,
		})
		return
	}
	helper.LogInfo("AMQP exchange declared", map[string]interface{}{
		"exchange_name":  exchangeName,
		"exchange_type":  exchangeType,
		"retry_exchange": exchangeName + RetrySuffix,
	})
	return
}

// Declare Queue
func (ch *amqpChannel) DeclareQueue(queue, queueType, exchangeName, routeKey string, interval int) (amqp.Queue, error) {
	arg := make(amqp.Table)
	arg["x-dead-letter-exchange"] = exchangeName + RetrySuffix
	if routeKey != "" {
		arg["x-dead-letter-routing-key"] = routeKey + RetrySuffix
	}

	arg["x-queue-type"] = queueType
	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		arg,   // arguments
	)
	if err != nil {
		helper.LogError("AMQP queue declare failed", err, queue, map[string]interface{}{
			"queue":       queue,
			"queue_type":  queueType,
			"exchange":    exchangeName,
			"routing_key": routeKey,
		})
		return q, err
	}
	// queue retry
	argretry := make(amqp.Table)
	argretry["x-dead-letter-exchange"] = exchangeName
	if routeKey != "" {
		argretry["x-dead-letter-routing-key"] = routeKey
	}
	argretry["x-message-ttl"] = interval
	argretry["x-queue-type"] = queueType
	retryQueueName := fmt.Sprintf("%s%s", queue, QueueRetrySuffix)
	_, err = ch.QueueDeclare(
		retryQueueName,
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		argretry, // arguments
	)
	if err != nil {
		helper.LogError("AMQP retry queue declare failed", err, retryQueueName, map[string]interface{}{
			"queue":      retryQueueName,
			"queue_type": queueType,
			"interval":   interval,
		})
		return q, err
	}
	helper.LogInfo("AMQP queue declared", map[string]interface{}{
		"queue":       queue,
		"queue_type":  queueType,
		"retry_queue": retryQueueName,
		"exchange":    exchangeName,
		"routing_key": routeKey,
		"interval":    interval,
	})
	return q, err
}

// Bind Queue
func (ch *amqpChannel) BindQueue(q amqp.Queue, routeKey, exchangeName string) error {
	var err error
	rKey := ""

	if routeKey != "" {
		rKey = routeKey + RetrySuffix
	}
	err = ch.QueueBind(
		q.Name,
		routeKey,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		helper.LogError("AMQP queue bind failed", err, q.Name, map[string]interface{}{
			"queue":       q.Name,
			"routing_key": routeKey,
			"exchange":    exchangeName,
		})
		return err
	}
	retryQueueName := fmt.Sprintf("%s%s", q.Name, QueueRetrySuffix)
	err = ch.QueueBind(
		retryQueueName,
		rKey,
		exchangeName+RetrySuffix,
		false,
		nil,
	)
	if err != nil {
		helper.LogError("AMQP retry queue bind failed", err, retryQueueName, map[string]interface{}{
			"queue":       retryQueueName,
			"routing_key": rKey,
			"exchange":    exchangeName + RetrySuffix,
		})
		return err
	}
	helper.LogInfo("AMQP queue bound", map[string]interface{}{
		"queue":       q.Name,
		"retry_queue": retryQueueName,
		"routing_key": routeKey,
		"exchange":    exchangeName,
	})
	return err
}

// New Queue
func (ch *amqpChannel) NewQueue(exchangeName, queueName, queueType, routeKey string, interval int) (amqp.Queue, error) {
	var err error
	var queue amqp.Queue
	if queue, err = ch.DeclareQueue(queueName, queueType, exchangeName, routeKey, interval); err != nil {
		helper.LogError("AMQP new queue failed at declare", err, queueName, map[string]interface{}{
			"queue":       queueName,
			"queue_type":  queueType,
			"exchange":    exchangeName,
			"routing_key": routeKey,
			"interval":    interval,
		})
		return queue, err
	}
	if err = ch.BindQueue(queue, routeKey, exchangeName); err != nil {
		helper.LogError("AMQP new queue failed at bind", err, queueName, map[string]interface{}{
			"queue":       queueName,
			"routing_key": routeKey,
			"exchange":    exchangeName,
		})
		return queue, err
	}
	helper.LogInfo("AMQP new queue created", map[string]interface{}{
		"queue":       queueName,
		"queue_type":  queueType,
		"exchange":    exchangeName,
		"routing_key": routeKey,
		"interval":    interval,
	})
	return queue, err
}

// Read Message
func (ch *amqpChannel) ReadMessage(q amqp.Queue) (<-chan amqp.Delivery, error) {
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err)
	return msgs, err
}

func (ch *amqpChannel) PublishMessage(q *amqp.Queue, routingKey, contentType, exchange string, body []byte) error {
	var key = ""

	if routingKey != "" {
		key = routingKey
	}

	if q != nil {
		key = q.Name
	}

	err := ch.Publish(
		exchange, // exchange
		key,      // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType:  contentType,
			Body:         body,
			DeliveryMode: 2,
		},
	)
	if err != nil {
		helper.LogError("AMQP publish message failed", err, exchange, map[string]interface{}{
			"exchange":     exchange,
			"routing_key":  key,
			"content_type": contentType,
			"body_size":    len(body),
		})
	} else {
		helper.LogInfo("AMQP message published", map[string]interface{}{
			"exchange":     exchange,
			"routing_key":  key,
			"content_type": contentType,
			"body_size":    len(body),
		})
	}
	return err
}

func failOnError(err error) {
	if err != nil {
		helper.LogError("AMQP operation failed", err, "", nil)
	}
}

func (ch *amqpChannel) GetChannel() *amqp.Channel {
	return ch.Channel
}
