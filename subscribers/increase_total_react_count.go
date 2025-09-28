package subscribers

import (
	"context"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/todo/storage/mysql"
	"github.com/coderconquerer/social-todo/pubsub"
	"gorm.io/gorm"
)

type TodoId interface {
	GetTodoId() int
}

// IncreaseTotalReactionCount returns a subJob that will increase the like count
// in the database when a user likes an item.
func IncreaseTotalReactionCount(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
		// Job title for identification/logging
		Title: "Increase total reaction count after user likes item",

		SHandler: func(ctx context.Context, message *pubsub.Message) error {
			// Get the main database connection from the service context
			db := serviceCtx.MustGet(common.DbMainName).(*gorm.DB)

			// Extract the item ID data from the incoming message
			data := message.Data().(TodoId)

			// Call the storage layer to increase the like count for the given item ID
			return mysql.GetNewMySQLConnection(db).IncreaseTotalReactionCount(ctx, data.GetTodoId())
		},
	}
}
