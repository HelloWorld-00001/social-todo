package subscribers

import (
	"context"
	"fmt"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/coderconquerer/social-todo/pubsub"
)

type TodoNotify interface {
	GetTodoId() int
	GetUserId() int
}

// IncreaseTotalReactionCount returns a subJob that will increase the like count
// in the database when a user likes an item.
func NotifyUserReactTodoItem(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
		// Job title for identification/logging
		Title: "Notify that an user has reacted a todo item",

		SHandler: func(ctx context.Context, message *pubsub.Message) error {

			data := message.Data().(TodoNotify)

			fmt.Printf("User with id %d has reacted item %d", data.GetUserId(), data.GetTodoId())
			// Call the storage layer to increase the like count for the given item ID
			return nil
		},
	}
}
