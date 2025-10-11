package storage

import (
	"context"
	"time"

	"github.com/coderconquerer/social-todo/grpc/contract"
)

type authRpcClient struct {
	client contract.AuthenticationServiceClient
}

// NewAuthClientGrpc connects to the gRPC server at given address (e.g. ":50051").
func NewAuthClientGrpc(client contract.AuthenticationServiceClient) *authRpcClient {
	return &authRpcClient{client: client}
}

// --------------------
// RPC Calls
// --------------------

// Login calls the Login RPC and returns the token.
func (c *authRpcClient) Login(ctx context.Context, username, password string) (string, error) {
	req := &contract.LoginRequest{
		Username: username,
		Password: password,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.Login(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.GetToken(), nil
}

func (c *authRpcClient) RegisterAccount(ctx context.Context, username, password string) (bool, error) {
	req := &contract.RegisterAccountRequest{
		Username: username,
		Password: password,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.RegisterAccount(ctx, req)
	if err != nil {
		return false, err
	}
	return resp.GetSuccess(), nil
}

func (c *authRpcClient) DisableAccount(ctx context.Context, id int32, disable int32) (bool, error) {
	req := &contract.DisableAccountRequest{
		Id:      id,
		Disable: disable,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.DisableAccount(ctx, req)
	if err != nil {
		return false, err
	}
	return resp.GetSuccess(), nil
}
