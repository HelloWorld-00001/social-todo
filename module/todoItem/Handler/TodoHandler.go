package Handler

import "github.com/coderconquerer/social-todo/module/todoItem/BusinessUseCases"

type TodoHandler struct {
	GetTodoDetailBz *BusinessUseCases.GetTodoDetailLogic
	GetTodoListBz   *BusinessUseCases.GetTodoListLogic
	CreateTodoBz    *BusinessUseCases.CreateTodoLogic
	DeleteTodoBz    *BusinessUseCases.DeleteTodoItemLogic
}

func NewTodoHandler(GetTodoDetail *BusinessUseCases.GetTodoDetailLogic,
	GetTodoList *BusinessUseCases.GetTodoListLogic,
	CreateTodo *BusinessUseCases.CreateTodoLogic,
	DeleteTodo *BusinessUseCases.DeleteTodoItemLogic) *TodoHandler {
	return &TodoHandler{
		GetTodoDetailBz: GetTodoDetail,
		GetTodoListBz:   GetTodoList,
		CreateTodoBz:    CreateTodo,
		DeleteTodoBz:    DeleteTodo,
	}
}
