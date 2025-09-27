package registerservice

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	url        string
	queue      amqp.Queue
	channel    *amqp.Channel
	connection *amqp.Connection
}

// NewRabbitMQConsumer initializes consumer with queue
func NewRabbitMQConsumer(url, queueName string) *RabbitMQConsumer {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}

	q, err := ch.QueueDeclare(
		queueName,
		false, // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	return &RabbitMQConsumer{
		url:        url,
		queue:      q,
		channel:    ch,
		connection: conn,
	}
}

// Close releases resources
func (c *RabbitMQConsumer) Close() {
	if err := c.channel.Close(); err != nil {
		log.Printf("Failed to close channel: %v", err)
	}
	if err := c.connection.Close(); err != nil {
		log.Printf("Failed to close connection: %v", err)
	}
}

// Consume starts consuming messages with a handler
func (c *RabbitMQConsumer) Consume(handler func(msg []byte)) error {
	msgs, err := c.channel.Consume(
		c.queue.Name,
		"",
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
		for d := range msgs {
			handler(d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages on queue %s. To exit press CTRL+C", c.queue.Name)
	select {} // block forever
}
