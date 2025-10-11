package business

import (
	"context"
	"errors"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/authenticationrpc/entity"
)

type AuthenticationGrpcService interface {
	Login(ctx context.Context, username, password string) (string, error)
	RegisterAccount(ctx context.Context, username, password string) (bool, error)
	DisableAccount(ctx context.Context, id int32, disable int32) (bool, error)
}

// AuthenticationBusinessGrpc defines the contract for authentication logic
type AuthenticationBusinessGrpc interface {
	Login(c context.Context, acc entity.AccountLogin) (string, error)
	RegisterAccount(c context.Context, acc *entity.AccountRegister) error
	DisableAccount(c context.Context, id, isDisable int) error
}

type authenticationBusinessGrpc struct {
	grpcService AuthenticationGrpcService
}

func NewAuthenticationBusinessGrpc(store AuthenticationGrpcService) AuthenticationBusinessGrpc {
	return &authenticationBusinessGrpc{grpcService: store}
}

func (bz *authenticationBusinessGrpc) Login(c context.Context, acc entity.AccountLogin) (string, error) {
	res, err := bz.grpcService.Login(c, acc.Username, acc.Password)

	if err != nil {
		return "", common.InternalServerError.WithError(errors.New("cannot login with current account")).WithRootCause(err)
	}

	return res, nil
}

func (bz *authenticationBusinessGrpc) RegisterAccount(c context.Context, acc *entity.AccountRegister) error {
	res, err := bz.grpcService.RegisterAccount(c, acc.Username, acc.Password)

	if err != nil || !res {
		return common.InternalServerError.WithError(entity.ErrCannotCreateAccount).WithRootCause(err)
	}
	return nil
}

func (bz *authenticationBusinessGrpc) DisableAccount(c context.Context, id, isDisable int) error {
	ok, err := bz.grpcService.DisableAccount(c, int32(id), int32(isDisable))
	if err != nil || !ok {
		return common.InternalServerError.WithError(entity.ErrCannotDisableAccount).WithRootCause(err)
	}

	return nil
}
