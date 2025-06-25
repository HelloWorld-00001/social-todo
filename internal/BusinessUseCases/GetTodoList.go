package BusinessUseCases

import (
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/coderconquerer/go-login-app/internal/models"
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

func (bz *GetTodoListLogic) GetTodoList(c *gin.Context, filter *common.Filter, pagination *common.Pagination) ([]models.Todo, error) {
	data, err := bz.store.GetTodoList(c, filter, pagination)
	if err != nil {
		return nil, err
	}
	return data, nil
}
