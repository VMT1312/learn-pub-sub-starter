package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	DurableQueue SimpleQueueType = iota
	TransientQueue
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	tbl := amqp.Table{
		"x-dead-letter-exchange": "peril_dlx",
	}

	var q amqp.Queue
	switch queueType {
	case DurableQueue:
		q, err = ch.QueueDeclare(
			queueName,
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			tbl,   // arguments
		)
	case TransientQueue:
		q, err = ch.QueueDeclare(
			queueName,
			false, // durable
			true,  // delete when unused
			true,  // exclusive
			false, // no-wait
			tbl,   // arguments
		)
	}
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return ch, q, nil
}
