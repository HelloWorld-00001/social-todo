package BusinessUseCases

import (
	"context"
	common "github.com/coderconquerer/social-todo/common"
	models2 "github.com/coderconquerer/social-todo/module/user/models"
	"github.com/gin-gonic/gin"
)

type GetReactedUsersStorage interface {
	GetReactedUsers(c *gin.Context, todoId int, pagination *common.Pagination) ([]models2.SimpleUser, error)
	GetReactedTodo(c context.Context, todoIds []int) (map[int]int, error)
}

type GetListReactedUsersLogic struct {
	store GetReactedUsersStorage
}

func GetNewGetListReactedUsersLogic(store GetReactedUsersStorage) *GetListReactedUsersLogic {
	return &GetListReactedUsersLogic{store: store}
}

func (bz *GetListReactedUsersLogic) GetListReactedUsers(c *gin.Context, todoId int, pagination *common.Pagination) ([]models2.SimpleUser, *common.AppError) {
	data, err := bz.store.GetReactedUsers(c, todoId, pagination)
	if err != nil {
		return nil, common.NewCannotGetEntity(models2.User{}.TableName(), err)
	}
	return data, nil
}

func (bz *GetListReactedUsersLogic) GetTodoItemTotalReact(c context.Context, todoIds []int) (map[int]int, error) {
	data, err := bz.store.GetReactedTodo(c, todoIds)
	if err != nil {
		return nil, err
	}
	return data, nil
}
