package BusinessUseCases

import (
	"github.com/coderconquerer/go-login-app/internal/TodoItem/models"
	"github.com/coderconquerer/go-login-app/internal/common"
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

func (bz *CreateTodoLogic) CreateTodoItem(c *gin.Context, todo *models.TodoCreation) error {
	if err := bz.store.CreateTodoItem(c, todo); err != nil {
		return common.NewCannotCreateEntity(models.TodoTableName, err)
	}
	return nil
}
