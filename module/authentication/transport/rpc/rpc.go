package rpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/coderconquerer/social-todo/common"

	"github.com/coderconquerer/social-todo/grpc/contract"
	"github.com/coderconquerer/social-todo/module/authentication/business"
	"github.com/coderconquerer/social-todo/module/authentication/entity"
)

type rpcServer struct {
	business business.AuthenticationBusiness
}

// NewRPCServer creates a new gRPC AuthenticationService server.
func NewRPCServer(b business.AuthenticationBusiness) contract.AuthenticationServiceServer {
	return &rpcServer{business: b}
}

// ---------------------------
// Implementations
// ---------------------------

func (s *rpcServer) Login(ctx context.Context, req *contract.LoginRequest) (*contract.LoginResponse, error) {
	fmt.Print("Received Login request from other service")
	acc := entity.AccountLogin{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}

	token, err := s.business.Login(ctx, acc)
	if err != nil {
		return nil, err
	}
	if token == nil {
		// Return empty token to indicate invalid login
		return &contract.LoginResponse{Token: ""}, nil
	}

	return &contract.LoginResponse{Token: token.GetToken()}, nil
}

func (s *rpcServer) RegisterAccount(ctx context.Context, req *contract.RegisterAccountRequest) (*contract.RegisterAccountResponse, error) {
	acc := entity.AccountRegister{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}

	if err := s.business.RegisterAccount(ctx, &acc); err != nil {
		return nil, err
	}

	return &contract.RegisterAccountResponse{Success: true}, nil
}

func (s *rpcServer) DisableAccount(ctx context.Context, req *contract.DisableAccountRequest) (*contract.DisableAccountResponse, error) {
	id := int(req.GetId())
	disable := req.GetDisable()

	if disable != 0 && disable != 1 {
		return nil, common.BadRequest.WithError(errors.New("disable must be zero or one"))
	}

	if err := s.business.DisableAccount(ctx, id, disable == 1); err != nil {
		return nil, err
	}

	return &contract.DisableAccountResponse{Success: true}, nil
}
