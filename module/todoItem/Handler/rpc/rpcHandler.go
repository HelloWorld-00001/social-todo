package rpc

import "github.com/coderconquerer/social-todo/module/todoItem/BusinessUseCases"

type RpcHandler struct {
	GetTodoItemTotalReactBz *BusinessUseCases.GetTodoItemTotalReactLogic
}

func NewRpcHandler(getTodoItemTotalReactBz *BusinessUseCases.GetTodoItemTotalReactLogic) *RpcHandler {
	return &RpcHandler{
		GetTodoItemTotalReactBz: getTodoItemTotalReactBz,
	}
}
