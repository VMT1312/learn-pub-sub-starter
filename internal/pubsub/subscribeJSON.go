package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int

const (
	Ack Acktype = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) Acktype,
) error {
	ch, q, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		defer ch.Close()
		for msg := range msgs {
			var v T
			err := json.Unmarshal(msg.Body, &v)
			if err != nil {
				msg.Ack(false)
				continue
			}
			ack := handler(v)
			switch ack {
			case Ack:
				log.Println("Acknowledging message")
				msg.Ack(false)
			case NackRequeue:
				log.Println("Negatively acknowledging message, requeueing")
				msg.Nack(false, true)
			case NackDiscard:
				log.Println("Negatively acknowledging message, discarding")
				msg.Nack(false, false)
			}
			msg.Ack(false)
		}
	}()
	return nil
}
