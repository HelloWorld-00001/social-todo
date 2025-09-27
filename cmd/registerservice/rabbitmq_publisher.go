package registerservice

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type RbPublisher interface {
	PublishMessage(ctx context.Context, body string) error
}

type RabbitMQPublisher struct {
	url        string
	queue      amqp.Queue
	channel    *amqp.Channel
	connection *amqp.Connection
}

// NewRabbitMQPublisher initializes the publisher and keeps resources open
func NewRabbitMQPublisher(url, queueName string) *RabbitMQPublisher {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	q, err := ch.QueueDeclare(
		queueName,
		false, // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	return &RabbitMQPublisher{
		url:        url,
		queue:      q,
		channel:    ch,
		connection: conn,
	}
}

// Close cleans up resources
func (rp *RabbitMQPublisher) Close() {
	if err := rp.channel.Close(); err != nil {
		log.Printf("Failed to close channel: %v", err)
	}
	if err := rp.connection.Close(); err != nil {
		log.Printf("Failed to close connection: %v", err)
	}
}

// PublishMessage publishes a message with context
func (rp *RabbitMQPublisher) PublishMessage(ctx context.Context, body string) error {
	err := rp.channel.PublishWithContext(ctx,
		"",            // exchange
		rp.queue.Name, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		return err
	}
	log.Printf(" [x] Sent %s\n", body)
	return nil
}
