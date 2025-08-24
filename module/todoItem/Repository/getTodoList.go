package Repository

import (
	"context"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/todoItem/models"
)

type TodoListStorage interface {
	GetTodoList(c context.Context, filter *common.Filter, pagination *common.Pagination) ([]models.Todo, error)
}

type ReactionCountService interface {
	GetTodoTotalReact(c context.Context, todoIds []int) (map[int]int, error)
}

type TodoListWithReactRepo struct {
	todoStore    TodoListStorage
	reactService ReactionCountService
}

func GetNewTodoListWithReactRepo(todoStore TodoListStorage, reactStore ReactionCountService) *TodoListWithReactRepo {
	return &TodoListWithReactRepo{
		todoStore:    todoStore,
		reactService: reactStore,
	}
}

func (r *TodoListWithReactRepo) GetTodoListWithReactCount(c context.Context, filter *common.Filter, pagination *common.Pagination) ([]models.Todo, error) {
	todos, err := r.todoStore.GetTodoList(c, filter, pagination)
	if err != nil {
		return nil, err
	}
	size := len(todos)

	// quick return
	if size <= 0 {
		return todos, nil
	}

	todoIds := make([]int, size)
	for i, t := range todos {
		todoIds[i] = t.Id
	}

	reactionMap, err := r.reactService.GetTodoTotalReact(c, todoIds)
	if err != nil {
		return nil, err
	}

	// Aggregate the reaction of to do
	for i := range todos {
		todos[i].TotalReact = reactionMap[todos[i].Id]
	}

	return todos, nil
}
