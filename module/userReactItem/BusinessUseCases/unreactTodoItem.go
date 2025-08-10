package BusinessUseCases

import (
	"errors"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/userReactItem/models"
	"github.com/coderconquerer/social-todo/pubsub"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

type DeleteReactionStorage interface {
	FindReaction(c *gin.Context, userId, todoId int) (*models.Reaction, error)
	DeleteReaction(c *gin.Context, userId, todoId int) error
}

type UnreactTodoItemLogic struct {
	store DeleteReactionStorage
	ps    pubsub.PubSub
}

func GetNewUnreactTodoItemLogic(store DeleteReactionStorage, ps pubsub.PubSub) *UnreactTodoItemLogic {
	return &UnreactTodoItemLogic{store, ps}
}

func (bz *UnreactTodoItemLogic) UnreactTodoItem(c *gin.Context, userId, todoId int) *common.AppError {
	react, err := bz.store.FindReaction(c, userId, todoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NewCannotGetEntity(models.Reaction{}.TableName(), err)
		}
		return common.NewDatabaseError(err)
	}

	err2 := bz.store.DeleteReaction(c, react.UserId, react.TodoId)
	if err2 != nil {
		return common.NewInternalSeverErrorResponse(err2, err2.Error(), err2.Error())
	}

	var rct models.Reaction
	rct.UserId = userId
	rct.TodoId = todoId

	err3 := bz.ps.Publish(c, common.TopicDecreaseTotalReact, pubsub.NewMessage(rct))

	if err3 != nil {
		log.Println(err3)
	}
	return nil
}
