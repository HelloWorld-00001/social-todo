package jwtProvider

import (
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/coderconquerer/go-login-app/internal/components/tokenProviders"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type jwtProvider struct {
	prefix    string
	secretKey string
}

type myClaims struct {
	jwt.RegisteredClaims
	MyPayload common.TokenPayload `json:"TokenPayload"`
}

type token struct {
	Token string
}

func GetNewJwtProvider(prefix, secretKey string) *jwtProvider {
	return &jwtProvider{
		prefix:    prefix,
		secretKey: secretKey,
	}
}

func (t *token) GetToken() string {
	return t.Token
}

func (j *jwtProvider) GenerateToken(data tokenProviders.TokenPayload, expiry int) (tokenProviders.Token, error) {
	if j.secretKey == "" {
		// todo: define a more specific message
		return nil, tokenProviders.ErrGenerateToken
	}

	claims := &myClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiry) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		MyPayload: common.TokenPayload{
			UserId: data.GetUserId(),
			Role:   data.GetRole(),
		},
	}

	genToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := genToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return nil, tokenProviders.ErrGenerateToken
	}

	return &token{Token: signedToken}, nil
}

func (j *jwtProvider) ValidateToken(tokenString string) (tokenProviders.TokenPayload, error) {
	if tokenString == "" {
		return nil, tokenProviders.ErrMissingToken
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&myClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.secretKey), nil
		},
	)

	if err != nil {
		return nil, tokenProviders.NewInValidTokenErr(err)
	}

	validator := *jwt.NewValidator()
	if err := validator.Validate(token.Claims.(*myClaims)); err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*myClaims); ok && token.Valid {
		return claims.MyPayload, nil
	}

	return nil, tokenProviders.ErrInvalidToken
}
