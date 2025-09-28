package rpc

import "github.com/coderconquerer/social-todo/module/todo/business"

type RpcHandler struct {
	GetTodoItemTotalReactBz *business.todoBusiness
}

func NewRpcHandler(getTodoItemTotalReactBz *business.todoBusiness) *RpcHandler {
	return &RpcHandler{
		GetTodoItemTotalReactBz: getTodoItemTotalReactBz,
	}
}
