package storage

import (
	"context"
	"time"

	"github.com/coderconquerer/social-todo/grpc/contract"
	"google.golang.org/grpc"
)

type AuthRpcClient struct {
	client contract.AuthenticationServiceClient
	conn   *grpc.ClientConn
}

// NewAuthClient connects to the gRPC server at given address (e.g. ":50051").
func NewAuthClient(conn *grpc.ClientConn) (*AuthRpcClient, error) {
	client := contract.NewAuthenticationServiceClient(conn)
	return &AuthRpcClient{client: client, conn: conn}, nil
}

// Close the connection when done
func (c *AuthRpcClient) Close() error {
	return c.conn.Close()
}

// --------------------
// RPC Calls
// --------------------

// Login calls the Login RPC and returns the token.
func (c *AuthRpcClient) Login(ctx context.Context, username, password string) (string, error) {
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

func (c *AuthRpcClient) RegisterAccount(ctx context.Context, username, password string) (bool, error) {
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

func (c *AuthRpcClient) DisableAccount(ctx context.Context, id int32, disable int32) (bool, error) {
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
