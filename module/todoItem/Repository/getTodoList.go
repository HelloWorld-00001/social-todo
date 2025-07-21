package Repository

import (
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/todoItem/models"
	"github.com/gin-gonic/gin"
)

type TodoListStorage interface {
	GetTodoList(c *gin.Context, filter *common.Filter, pagination *common.Pagination) ([]models.Todo, error)
}

type ReactionCountStorage interface {
	GetReactedTodo(c *gin.Context, todoIds []int) (map[int]int, error)
}

type TodoListWithReactRepo struct {
	todoStore  TodoListStorage
	reactStore ReactionCountStorage
}

func GetNewTodoListWithReactRepo(todoStore TodoListStorage, reactStore ReactionCountStorage) *TodoListWithReactRepo {
	return &TodoListWithReactRepo{
		todoStore:  todoStore,
		reactStore: reactStore,
	}
}

func (r *TodoListWithReactRepo) GetTodoListWithReactCount(c *gin.Context, filter *common.Filter, pagination *common.Pagination) ([]models.Todo, error) {
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

	reactionMap, err := r.reactStore.GetReactedTodo(c, todoIds)
	if err != nil {
		return nil, err
	}

	// Aggregate the reaction of to do
	for i := range todos {
		todos[i].TotalReact = reactionMap[todos[i].Id]
	}

	return todos, nil
}
