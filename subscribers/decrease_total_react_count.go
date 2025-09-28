package subscribers

import (
	"context"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/todo/storage/mysql"
	"github.com/coderconquerer/social-todo/pubsub"
	"gorm.io/gorm"
)

// DecreaseTotalReactionCount returns a subJob that will decrease the like count
// in the database when a user likes an item.
func DecreaseTotalReactionCount(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
		// Job title for identification/logging
		Title: "Decrease total reaction count after user likes item",

		SHandler: func(ctx context.Context, message *pubsub.Message) error {
			// Get the main database connection from the service context
			db := serviceCtx.MustGet(common.DbMainName).(*gorm.DB)

			// Extract the item ID data from the incoming message
			data := message.Data().(TodoId)

			// Call the storage layer to increase the like count for the given item ID
			return mysql.GetNewMySQLConnection(db).DecreaseTotalReactionCount(ctx, data.GetTodoId())
		},
	}
}
