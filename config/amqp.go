package config

import (
	"boilerblade/config/amqp"
	"boilerblade/helper"
	"fmt"
)

func (e *Env) InitAMQP() *amqp.IAMQPConnection {
	amqpConn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", e.AMQP_USER, e.AMQP_PASSWORD, e.AMQP_HOST, e.AMQP_PORT))
	if err != nil {
		helper.LogError("AMQP connection failed", err, e.AMQP_HOST, map[string]interface{}{
			"host":     e.AMQP_HOST,
			"port":     e.AMQP_PORT,
			"user":     e.AMQP_USER,
			"password": "***",
			"url":      fmt.Sprintf("amqp://%s:%s@%s:%s/", e.AMQP_USER, e.AMQP_PASSWORD, e.AMQP_HOST, e.AMQP_PORT),
		})
		return nil
	}
	return &amqpConn
}
