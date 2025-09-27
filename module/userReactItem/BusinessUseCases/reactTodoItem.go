package BusinessUseCases

import (
	"context"
	"errors"
	"github.com/coderconquerer/social-todo/cmd/registerservice"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/userReactItem/models"
	"github.com/coderconquerer/social-todo/pubsub"
	"gorm.io/gorm"
	"log"
)

type CreateReactionStorage interface {
	CreateReaction(c context.Context, reaction models.Reaction) error
}

type ReactTodoItemLogic struct {
	store    CreateReactionStorage
	ps       pubsub.PubSub
	rabbitPs registerservice.RbPublisher
}

func GetNewReactTodoItemLogic(store CreateReactionStorage, ps pubsub.PubSub, rabbitPs registerservice.RbPublisher) *ReactTodoItemLogic {
	return &ReactTodoItemLogic{store, ps, rabbitPs}
}

func (bz *ReactTodoItemLogic) ReactTodoItem(c context.Context, reaction models.Reaction) *common.AppError {
	err := bz.store.CreateReaction(c, reaction)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NewCannotGetEntity(models.Reaction{}.TableName(), err)
		}
		return common.NewDatabaseError(err)
	}

	err2 := bz.ps.Publish(c, common.TopicIncreaseTotalReact, pubsub.NewMessage(reaction))
	if err2 != nil {
		log.Println(err2)
	}
	return nil
}
