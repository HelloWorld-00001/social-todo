package repository

import (
	"context"
	"fmt"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/todo/entity"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type TodoListStorage interface {
	GetTodoList(c context.Context, filter *common.Filter, pagination *common.Pagination) ([]entity.Todo, error)
}

type ReactionCountService interface {
	GetTodoTotalReact(c context.Context, todoIds []int) (map[int]int, error)
}

type TodoListWithReactRepo struct {
	todoStore    TodoListStorage
	reactService ReactionCountService
}

func GetNewTodoListWithReactRepo(todoStore TodoListStorage, reactStore ReactionCountService) *TodoListWithReactRepo {
	return &TodoListWithReactRepo{
		todoStore:    todoStore,
		reactService: reactStore,
	}
}

func (r *TodoListWithReactRepo) GetTodoListWithReactCount(c context.Context, filter *common.Filter, pagination *common.Pagination) ([]entity.Todo, error) {
	todos, err := r.todoStore.GetTodoList(c, filter, pagination)
	if err != nil {
		return nil, err
	}
	size := len(todos)

	// quick return
	if size <= 0 {
		return todos, nil
	}

	todoIds := make([]int, size)
	for i, t := range todos {
		todoIds[i] = t.Id
	}

	go func() {
		fmt.Printf(" [x] Requesting fib(%d)", 10)
		res, err := fibonacciRPC(10)
		fmt.Println(err, "Failed to handle RPC request")
		fmt.Printf(" [.] Got %d", res)
	}()

	reactionMap, err := r.reactService.GetTodoTotalReact(c, todoIds)
	if err != nil {
		return nil, err
	}
	// Aggregate the reaction of to do
	for i := range todos {
		todos[i].TotalReact = reactionMap[todos[i].Id]
	}

	return todos, nil
}

func fibonacciRPC(n int) (res int, err error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	log.Printf("Failed to connect to RabbitMQ server: %v", err)
	fmt.Printf("Failed to connect to RabbitMQ server: %v", err)

	defer conn.Close()

	ch, err := conn.Channel()
	fmt.Print(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	fmt.Print(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	fmt.Print(err, "Failed to register a consumer")

	corrId := randomString(32)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(strconv.Itoa(n)),
		})
	fmt.Print(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			res, err = strconv.Atoi(string(d.Body))
			fmt.Print(err, "Failed to convert body to integer")
			break
		}
	}

	return
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
