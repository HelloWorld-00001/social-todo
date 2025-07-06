package BusinessUseCases

import (
	"github.com/coderconquerer/go-login-app/internal/TodoItem/models"
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/gin-gonic/gin"
)

type GetTodoListStorage interface {
	GetTodoList(c *gin.Context, filter *common.Filter, pagination *common.Pagination) ([]models.Todo, error)
}

type GetTodoListLogic struct {
	store GetTodoListStorage
}

func GetNewGetTodoListLogic(store GetTodoListStorage) *GetTodoListLogic {
	return &GetTodoListLogic{store: store}
}

func (bz *GetTodoListLogic) GetTodoList(c *gin.Context, filter *common.Filter, pagination *common.Pagination) ([]models.Todo, *common.AppError) {
	data, err := bz.store.GetTodoList(c, filter, pagination)
	if err != nil {
		return nil, common.NewCannotGetEntity(models.TodoTableName, err)
	}
	return data, nil
}
