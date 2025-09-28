package rpc

import "github.com/coderconquerer/social-todo/module/todo/business"

type RpcHandler struct {
	GetTodoItemTotalReactBz business.TodoBusiness
}

func NewRpcHandler(getTodoItemTotalReactBz business.TodoBusiness) *RpcHandler {
	return &RpcHandler{
		GetTodoItemTotalReactBz: getTodoItemTotalReactBz,
	}
}
