package BusinessUseCases

import (
	"errors"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/userReactItem/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateReactionStorage interface {
	CreateReaction(c *gin.Context, reaction models.Reaction) error
}

type ReactTodoItemLogic struct {
	store CreateReactionStorage
}

func GetNewReactTodoItemLogic(store CreateReactionStorage) *ReactTodoItemLogic {
	return &ReactTodoItemLogic{store: store}
}

func (bz *ReactTodoItemLogic) ReactTodoItem(c *gin.Context, reaction models.Reaction) *common.AppError {
	err := bz.store.CreateReaction(c, reaction)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NewCannotGetEntity(models.Reaction{}.TableName(), err)
		}
		return common.NewDatabaseError(err)
	}

	return nil
}
