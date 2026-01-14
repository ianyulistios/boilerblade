package amqp

import (
	"boilerblade/helper"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

type connection struct {
	*amqp.Connection
}

// Dial wrap amqp.Dial, dial and get a reconnect connection
func Dial(url string) (IAMQPConnection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	connection := &connection{
		Connection: conn,
	}

	go func() {
		for {
			reason, ok := <-connection.Connection.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok {
				helper.LogInfo("AMQP connection closed by developer", map[string]interface{}{
					"url": url,
				})
				break
			}
			helper.LogError("AMQP connection closed", fmt.Errorf(reason.Error()), url, map[string]interface{}{
				"reason": reason.Reason,
				"code":   reason.Code,
			})

			// reconnect if not closed by developer
			for {
				// wait before reconnect
				time.Sleep(delay * time.Second)

				conn, err := amqp.Dial(url)
				if err == nil {
					connection.Connection = conn
					helper.LogInfo("AMQP reconnect success", map[string]interface{}{
						"url": url,
					})
					break
				}

				helper.LogError("AMQP reconnect failed", err, url, nil)
			}
		}
	}()

	return connection, nil
}

// Channel wrap amqp.Connection.Channel, get a auto reconnect channel
func (c *connection) Channel() (IAMQPChannel, error) {
	ch, err := c.Connection.Channel()
	if err != nil {
		return nil, err
	}

	prefetchCount := 20
	setChannelQoS(ch, prefetchCount)

	channel := &amqpChannel{Channel: ch}

	go func() {
		for {
			reason, ok := <-channel.Channel.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok || channel.IsClosed() {
				helper.LogInfo("AMQP channel closed", map[string]interface{}{})
				channel.Close() // close again, ensure closed flag set when connection closed
				break
			}
			helper.LogError("AMQP channel closed", fmt.Errorf(reason.Error()), "", map[string]interface{}{
				"reason": reason.Reason,
				"code":   reason.Code,
			})

			// reconnect if not closed by developer
			for {
				// wait for connection reconnect
				time.Sleep(delay * time.Second)

				ch, err := c.Connection.Channel()
				if err == nil {
					// Apply QoS settings to the recreated channel
					if err := setChannelQoS(ch, prefetchCount); err != nil {
						helper.LogError("AMQP channel QoS setting failed", err, "", nil)
					}
					helper.LogInfo("AMQP channel recreate success", map[string]interface{}{})
					channel.Channel = ch
					break
				}

				helper.LogError("AMQP channel recreate failed", err, "", nil)
			}
		}

	}()

	return channel, nil
}

func (c *connection) Close() error {
	return c.Connection.Close()
}

// setChannelQoS sets the Quality of Service settings for a channel
func setChannelQoS(ch *amqp.Channel, prefetchCount int) error {
	if err := ch.Qos(
		prefetchCount, // Prefetch count
		0,             // Prefetch size (0 means no limit)
		false,         // Global (applies to the current channel only)
	); err != nil {
		helper.LogError("AMQP channel QoS setting failed", err, "", map[string]interface{}{
			"prefetch_count": prefetchCount,
		})
		return err
	}
	return nil
}
