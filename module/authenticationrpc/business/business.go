package business

import (
	"context"
	"errors"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/grpc/contract"
	"github.com/coderconquerer/social-todo/module/authenticationrpc/entity"
	tokenProviders "github.com/coderconquerer/social-todo/plugin/tokenProviders"
)

type AuthenticationGrpcService interface {
	contract.AuthenticationServiceClient
}

// AuthenticationBusinessGrpc defines the contract for authentication logic
type AuthenticationBusinessGrpc interface {
	Login(c context.Context, acc entity.AccountLogin) (tokenProviders.Token, error)
	RegisterAccount(c context.Context, acc *entity.AccountRegister) error
	DisableAccount(c context.Context, id int, isDisable bool) error
}

type authenticationBusinessGrpc struct {
	grpcService   AuthenticationGrpcService
	tokenProvider tokenProviders.TokenProvider
	expiry        int
}

func NewAuthenticationBusinessGrpc(store AuthenticationGrpcService, tokenProvider tokenProviders.TokenProvider, expiry int) AuthenticationBusinessGrpc {
	return &authenticationBusinessGrpc{grpcService: store, tokenProvider: tokenProvider, expiry: expiry}
}

func (bz *authenticationBusinessGrpc) Login(c context.Context, acc entity.AccountLogin) (tokenProviders.Token, error) {

	res, err := bz.grpcService.Login(c, &contract.LoginRequest{Username: acc.Username, Password: acc.Password})

	if err != nil {
		return nil, common.InternalServerError.WithError(errors.New("cannot login with current account")).WithRootCause(err)
	}

	return res, nil
}

func (bz *authenticationBusinessGrpc) RegisterAccount(c context.Context, acc *entity.AccountRegister) error {
	res, err := bz.grpcService.RegisterAccount(c, &contract.RegisterAccountRequest{Username: acc.Username, Password: acc.Password})

	if err != nil || !res.Success {
		return common.InternalServerError.WithError(entity.ErrCannotCreateAccount).WithRootCause(err)
	}
	return nil
}

func (bz *authenticationBusinessGrpc) DisableAccount(c context.Context, id int, isDisable bool) error {

	res, err := bz.grpcService.DisableAccount(c, &contract.DisableAccountRequest{Id: int32(id)})
	if err != nil || !res.Success {
		return common.InternalServerError.WithError(entity.ErrCannotDisableAccount).WithRootCause(err)
	}

	return nil
}
