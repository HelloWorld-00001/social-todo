package BusinessUseCases

import (
	"context"
	common2 "github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/todoItem/models"
)

type GetTodoListRepo interface {
	GetTodoListWithReactCount(c context.Context, filter *common2.Filter, pagination *common2.Pagination) ([]models.Todo, error)
}

type GetTodoListLogic struct {
	repo GetTodoListRepo
}

func GetNewGetTodoListLogic(repo GetTodoListRepo) *GetTodoListLogic {
	return &GetTodoListLogic{repo: repo}
}

func (bz *GetTodoListLogic) GetTodoList(c context.Context, filter *common2.Filter, pagination *common2.Pagination) ([]models.Todo, *common2.AppError) {
	data, err := bz.repo.GetTodoListWithReactCount(c, filter, pagination)
	if err != nil {
		return nil, common2.NewCannotGetEntity(models.TodoTableName, err)
	}
	return data, nil
}
