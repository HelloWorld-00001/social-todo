package entity

import "errors"

var (
	ErrCannotCreateTodo    = errors.New("cannot create todo item")
	ErrCannotDeleteTodo    = errors.New("cannot delete todo item")
	ErrCannotGetTodo       = errors.New("cannot get todo item detail")
	ErrCannotGetList       = errors.New("cannot get todo list")
	ErrCannotGetReact      = errors.New("cannot get todo react count")
	ErrInvalidTodoCreation = errors.New("invalid todo creation parameter")
	ErrInvalidTodoUpdate   = errors.New("invalid todo update parameter")
)
