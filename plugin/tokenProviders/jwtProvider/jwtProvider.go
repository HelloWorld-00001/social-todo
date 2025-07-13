package jwtProvider

import (
	"flag"
	"github.com/coderconquerer/go-login-app/common"
	tokenProviders "github.com/coderconquerer/go-login-app/plugin/tokenProviders"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtProvider struct {
	prefix    string
	secretKey string
}

func (j *JwtProvider) GetPrefix() string {
	return j.prefix
}

func (j *JwtProvider) Get() interface{} {
	// Return itself so it can be used dynamically (via interface{})
	return j
}

func (j *JwtProvider) Name() string {
	// Unique name of this plugin
	return "jwt-provider"
}

func (j *JwtProvider) InitFlags() {
	// Define command-line flags or config overrides here
	flag.StringVar(&j.secretKey, "jwt-secret-key", "", "Secret key used for signing JWTs")
	flag.StringVar(&j.prefix, "jwt-prefix", "jwt", "S3Prefix for the JWT token")
}

func (j *JwtProvider) Configure() error {
	// Validate essential configuration
	if j.secretKey == "" {
		return tokenProviders.ErrGenerateToken
	}
	if j.prefix == "" {
		j.prefix = "jwt-prefix"
	}
	return nil
}

func (j *JwtProvider) Run() error {
	// No background task is needed for JWT
	return nil
}

func (j *JwtProvider) Stop() <-chan bool {
	// Return a closed channel indicating nothing to clean up
	ch := make(chan bool)
	go func() {
		ch <- true
	}()
	return ch
}

type myClaims struct {
	jwt.RegisteredClaims
	MyPayload common.TokenPayload `json:"TokenPayload"`
}

type token struct {
	Token string
}

func GetNewJwtProvider(prefix, secretKey string) *JwtProvider {
	return &JwtProvider{
		prefix:    prefix,
		secretKey: secretKey,
	}
}

func (t *token) GetToken() string {
	return t.Token
}

func (j *JwtProvider) GenerateToken(data tokenProviders.TokenPayload, expiry int) (tokenProviders.Token, error) {
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

func (j *JwtProvider) ValidateToken(tokenString string) (tokenProviders.TokenPayload, error) {
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
