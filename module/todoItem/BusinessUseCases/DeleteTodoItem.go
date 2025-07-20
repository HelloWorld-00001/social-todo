package BusinessUseCases

import (
	"errors"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/todoItem/models"
	"github.com/gin-gonic/gin"
)

type DeleteTodoItemStorage interface {
	DeleteTodoItem(c *gin.Context, id int) error
}

type DeleteTodoItemLogic struct {
	store DeleteTodoItemStorage
}

func GetNewDeleteTodoItemLogic(store DeleteTodoItemStorage) *DeleteTodoItemLogic {
	return &DeleteTodoItemLogic{store: store}
}

func (bz *DeleteTodoItemLogic) DeleteTodoItem(c *gin.Context, id int) *common.AppError {
	err := bz.store.DeleteTodoItem(c, id)
	if !errors.Is(err, nil) {
		return common.NewCannotDeleteEntity(models.TodoTableName, err)
	}
	return nil
}
