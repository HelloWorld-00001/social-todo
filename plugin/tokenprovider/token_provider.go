package plugin

import (
	"errors"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/plugin"
)

type TokenProvider interface {
	GenerateToken(data TokenPayload, expiry int) (Token, error)
	ValidateToken(token string) (TokenPayload, error)
	plugin.PluginBase
}

type TokenPayload interface {
	GetAccountId() int
	GetUserId() int
	GetRole() string
}

type Token interface {
	GetToken() string
}

func NewInValidTokenErr(err error) error {
	return common.BadRequest.WithError(ErrInvalidToken).WithRootCause(err)
}

var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrExpiredToken   = errors.New("token has expired")
	ErrMalformedToken = errors.New("token is malformed")
	ErrGenerateToken  = errors.New("cannot generate token")
	ErrMissingToken   = errors.New("token is missing from request")
)
