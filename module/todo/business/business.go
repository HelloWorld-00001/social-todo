package business

import (
	"context"
	"errors"

	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/todo/entity"
)

// TodoStorage Bundle of storage/repo dependencies
type TodoStorage interface {
	CreateTodoItem(c context.Context, todo *entity.TodoCreation) error
	DeleteTodoItem(c context.Context, id int) error
	GetTodoItemDetailById(c context.Context, id int) (*entity.Todo, error)
}

type TodoListWithReactRepo interface {
	//GetReactedTodo(c context.Context, todoIds []int) (map[int]int, error)
	GetTodoListWithReactCount(c context.Context, filter *common.Filter, pagination *common.Pagination) ([]entity.Todo, error)
}

// TodoBusiness defines the contract for todo business logic.
type TodoBusiness interface {
	CreateTodoItem(c context.Context, todo *entity.TodoCreation) error
	DeleteTodoItem(c context.Context, id int) error
	GetTodoDetail(c context.Context, id int) (*entity.Todo, error)
	GetTodoList(c context.Context, filter *common.Filter, pagination *common.Pagination) ([]entity.Todo, error)
	//GetTodoItemTotalReact(c context.Context, todoIds []int) (map[int]int, error)
}

// todoBusiness Main business struct
type todoBusiness struct {
	store             TodoStorage
	todoWithReactRepo TodoListWithReactRepo
}

// NewTodoBusiness Constructor
func NewTodoBusiness(store TodoStorage, todoWithReactRepo TodoListWithReactRepo) TodoBusiness {
	return &todoBusiness{store, todoWithReactRepo}
}

// ========== Actions ==========

// CreateTodoItem Create new todoTask
func (bz *todoBusiness) CreateTodoItem(c context.Context, todo *entity.TodoCreation) error {
	if err := bz.store.CreateTodoItem(c, todo); !errors.Is(err, nil) {
		return common.InternalServerError.WithError(entity.ErrCannotCreateTodo).WithRootCause(err)
	}
	return nil
}

// DeleteTodoItem Delete a todoTask
func (bz *todoBusiness) DeleteTodoItem(c context.Context, id int) error {
	err := bz.store.DeleteTodoItem(c, id)
	if !errors.Is(err, nil) {
		return common.InternalServerError.WithError(entity.ErrCannotDeleteTodo).WithRootCause(err)
	}
	return nil
}

// GetTodoDetail Get a detail todoTask
func (bz *todoBusiness) GetTodoDetail(c context.Context, id int) (*entity.Todo, error) {
	data, err := bz.store.GetTodoItemDetailById(c, id)
	if err != nil {
		return nil, common.InternalServerError.WithError(entity.ErrCannotGetTodo).WithRootCause(err)
	}
	return data, nil
}

// GetTodoItemTotalReact get react count
//func (bz *todoBusiness) GetTodoItemTotalReact(c context.Context, todoIds []int) (map[int]int, error) {
//	data, err := bz.todoWithReactRepo.GetReactedTodo(c, todoIds)
//	if err != nil {
//		return nil, common.InternalServerError.WithError(entity.ErrCannotGetReact).WithRootCause(err)
//	}
//	return data, nil
//}

// GetTodoList get todoTask list
func (bz *todoBusiness) GetTodoList(c context.Context, filter *common.Filter, pagination *common.Pagination) ([]entity.Todo, error) {
	data, err := bz.todoWithReactRepo.GetTodoListWithReactCount(c, filter, pagination)
	if err != nil {
		return nil, common.InternalServerError.WithError(entity.ErrCannotGetList).WithRootCause(err)
	}
	return data, nil
}
