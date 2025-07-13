package Handler

import "github.com/coderconquerer/go-login-app/module/account/BusinessUseCases"

type AccountHandler struct {
	LoginLogic    *BusinessUseCases.LoginLogic
	RegisterLogic *BusinessUseCases.RegisterAccountLogic
	DisableLogic  *BusinessUseCases.DisableAccountLogic
}

func NewAccountHandler(LoginLogic *BusinessUseCases.LoginLogic, RegisterLogic *BusinessUseCases.RegisterAccountLogic, DisableLogic *BusinessUseCases.DisableAccountLogic) *AccountHandler {
	return &AccountHandler{
		LoginLogic:    LoginLogic,
		RegisterLogic: RegisterLogic,
		DisableLogic:  DisableLogic,
	}
}
