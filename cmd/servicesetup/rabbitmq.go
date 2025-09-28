package servicesetup

import (
	"context"
	"go.uber.org/zap"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartRabbitMQConsumer(logger *zap.Logger) {
	go processFib(logger)
}

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}

func processFib(logger *zap.Logger) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logger.Warn("Failed to connect to RabbitMQ: %v", zap.Error(err))
		return
	}
	defer conn.Close()

	ch, errChan := conn.Channel()
	if errChan != nil {
		logger.Warn("Failed to open a channel: %v", zap.Error(errChan))
		return
	}
	defer ch.Close()

	q, errQueue := ch.QueueDeclare("rpc_queue", false, false, false, false, nil)
	if errQueue != nil {
		logger.Warn("Failed to declare a queue: %v", zap.Error(errQueue))
		return
	}

	if errQos := ch.Qos(1, 0, false); errQos != nil {
		logger.Warn("Failed to set QoS: %v", zap.Error(errQos))
		return
	}

	msgs, errConsume := ch.Consume(q.Name, "", false, false, false, false, nil)
	if errConsume != nil {
		logger.Warn("Failed to register a consumer: %v", zap.Error(errConsume))
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for d := range msgs {
			n, _ := strconv.Atoi(string(d.Body))
			response := fib(n)

			err = ch.PublishWithContext(ctx, "",
				d.ReplyTo, false, false,
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(strconv.Itoa(response)),
				})
			if err != nil {
				logger.Warn("Failed to publish a message: %v", zap.Error(err))
			}
			d.Ack(false)
		}
	}()

	logger.Info("[*] Awaiting RPC requests")
	select {}
}
