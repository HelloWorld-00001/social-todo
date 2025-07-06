package BusinessUseCases

import (
	"github.com/coderconquerer/go-login-app/internal/TodoItem/models"
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/gin-gonic/gin"
)

type GetTodoDetailStorage interface {
	GetTodoItemDetailById(c *gin.Context, id int) (*models.Todo, error)
}

type GetTodoDetailLogic struct {
	store GetTodoDetailStorage
}

func GetNewGetTodoDetailLogic(store GetTodoDetailStorage) *GetTodoDetailLogic {
	return &GetTodoDetailLogic{store: store}
}

func (bz *GetTodoDetailLogic) GetTodoDetail(c *gin.Context, id int) (*models.Todo, *common.AppError) {
	data, err := bz.store.GetTodoItemDetailById(c, id)
	if err != nil {
		return nil, common.NewCannotGetEntity(models.TodoTableName, err)
	}
	return data, nil
}
