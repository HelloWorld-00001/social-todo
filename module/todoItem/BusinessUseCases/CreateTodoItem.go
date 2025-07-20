package BusinessUseCases

import (
	"errors"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/todoItem/models"
	"github.com/gin-gonic/gin"
)

type CreateTodoStorage interface {
	CreateTodoItem(c *gin.Context, todo *models.TodoCreation) error
}

type CreateTodoLogic struct {
	store CreateTodoStorage
}

func GetNewCreateTodoLogic(store CreateTodoStorage) *CreateTodoLogic {
	return &CreateTodoLogic{store: store}
}

func (bz *CreateTodoLogic) CreateTodoItem(c *gin.Context, todo *models.TodoCreation) *common.AppError {
	if err := bz.store.CreateTodoItem(c, todo); !errors.Is(err, nil) {
		return common.NewCannotCreateEntity(models.TodoTableName, err)
	}
	return nil
}
