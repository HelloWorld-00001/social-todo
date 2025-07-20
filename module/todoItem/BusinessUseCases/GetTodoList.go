package BusinessUseCases

import (
	common2 "github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/todoItem/models"
	"github.com/gin-gonic/gin"
)

type GetTodoListStorage interface {
	GetTodoList(c *gin.Context, filter *common2.Filter, pagination *common2.Pagination) ([]models.Todo, error)
}

type GetTodoListLogic struct {
	store GetTodoListStorage
}

func GetNewGetTodoListLogic(store GetTodoListStorage) *GetTodoListLogic {
	return &GetTodoListLogic{store: store}
}

func (bz *GetTodoListLogic) GetTodoList(c *gin.Context, filter *common2.Filter, pagination *common2.Pagination) ([]models.Todo, *common2.AppError) {
	data, err := bz.store.GetTodoList(c, filter, pagination)
	if err != nil {
		return nil, common2.NewCannotGetEntity(models.TodoTableName, err)
	}
	return data, nil
}
