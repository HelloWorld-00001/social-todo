package Handler

import "github.com/coderconquerer/go-login-app/module/user/BusinessUseCases"

type UserHandler struct {
	GetUserBz *BusinessUseCases.FindUserLogic
}

func NewUserHandler(GetUserBz *BusinessUseCases.FindUserLogic) *UserHandler {
	return &UserHandler{
		GetUserBz: GetUserBz,
	}
}
