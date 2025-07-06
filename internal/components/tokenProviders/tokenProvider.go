package tokenProviders

import (
	"errors"
	"github.com/coderconquerer/go-login-app/internal/common"
)

type TokenProvider interface {
	GenerateToken(data TokenPayload, expiry int) (Token, error)
	ValidateToken(token string) (TokenPayload, error)
}

type TokenPayload interface {
	GetUserId() int
	GetRole() string
}

type Token interface {
	GetToken() string
}

func NewInValidTokenErr(err error) *common.AppError {
	return common.NewBadRequestErrorResponse(
		err,
		"The token provided is invalid",
		err.Error(),
	)
}

var (
	ErrInvalidToken = common.NewCustomErrorResponse(
		errors.New("token is invalid"),
		"The token provided is invalid",
		"ErrInvalidToken",
	)

	ErrExpiredToken = common.NewCustomErrorResponse(
		errors.New("token has expired"),
		"The token has expired",
		"ErrExpiredToken",
	)

	ErrMalformedToken = common.NewCustomErrorResponse(
		errors.New("token is malformed"),
		"The token is malformed or corrupted",
		"ErrMalformedToken",
	)

	ErrGenerateToken = common.NewCustomErrorResponse(
		errors.New("cannot generate token"),
		"Unable to generate a new token",
		"ErrGenerateToken",
	)

	ErrMissingToken = common.NewCustomErrorResponse(
		errors.New("token is missing from request"),
		"No token provided in request",
		"ErrMissingToken",
	)
)
