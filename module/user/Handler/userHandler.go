package Handler

import "github.com/coderconquerer/social-todo/module/user/BusinessUseCases"

type UserHandler struct {
	GetUserBz *BusinessUseCases.FindUserLogic
}

func NewUserHandler(GetUserBz *BusinessUseCases.FindUserLogic) *UserHandler {
	return &UserHandler{
		GetUserBz: GetUserBz,
	}
}
