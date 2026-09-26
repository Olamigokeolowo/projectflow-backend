package events

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const queueName = "decision_created"

type RabbitMQQueue struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQQueue(url string) (*RabbitMQQueue, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable — survives a RabbitMQ restart
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQQueue{conn: conn, channel: ch}, nil
}

func (q *RabbitMQQueue) Publish(ctx context.Context, event DecisionCreated) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return q.channel.PublishWithContext(ctx,
		"",        // default exchange
		queueName, // routing key = queue name
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // survives a RabbitMQ restart, same as your durable queue
		},
	)
}

func (q *RabbitMQQueue) Close() {
	q.channel.Close()
	q.conn.Close()
}

func StartRabbitMQWorker(ctx context.Context, q *RabbitMQQueue) error {
	msgs, err := q.channel.Consume(
		queueName,
		"",    // consumer tag
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				var event DecisionCreated
				if err := json.Unmarshal(msg.Body, &event); err != nil {
					log.Println("failed to unmarshal event:", err)
					continue
				}
				log.Printf("processing event: decision %s created by %s (title: %s)",
					event.DecisionID, event.OwnerID, event.Title)
			case <-ctx.Done():
				log.Println("RabbitMQ worker shutting down")
				return
			}
		}
	}()

	return nil
}