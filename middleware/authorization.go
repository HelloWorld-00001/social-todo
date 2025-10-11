package middleware

import (
	"context"
	"errors"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/authentication/entity"
	tokenProviders "github.com/coderconquerer/social-todo/plugin/tokenprovider"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type AuthorizationStore interface {
	FindAccount(c context.Context, conditions map[string]interface{}) (*entity.Account, error)
}

func ErrorWrongAuthHeader(err error) error {
	return common.Unauthorized.WithError(errors.New("wrong authorization header")).WithRootCause(err)
}

func extractTokenHeader(authHeader string) (string, error) {
	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" || strings.TrimSpace(parts[1]) == "" {
		return "", ErrorWrongAuthHeader(errors.New("invalid authorization header"))
	}
	return parts[1], nil
}

// RequireAuth returns a Gin middleware function that checks for a valid token
func RequireAuth(provider tokenProviders.TokenProvider, store AuthorizationStore, roles ...string) func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString, err := extractTokenHeader(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorWrongAuthHeader(err))
			return
		}

		payload, errValidate := provider.ValidateToken(tokenString)
		if errValidate != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.Unauthorized.WithError(errors.New("invalid or expired token")).WithRootCause(errValidate))
			return
		}

		condition := map[string]interface{}{
			"Id": payload.GetAccountId(),
		}

		ctx := c.Request.Context()
		account, errFindAccount := store.FindAccount(ctx, condition)
		if err != nil || account == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.Unauthorized.WithError(errors.New("account not found")).WithRootCause(errFindAccount))
			return
		}

		if account.User == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.Unauthorized.WithError(errors.New("cannot get user information")))
			return
		}

		userRole := payload.GetRole()
		if len(roles) > 0 {
			allowed := false
			for _, role := range roles {
				if userRole == role {
					allowed = true
					break
				}
			}
			if !allowed {
				c.AbortWithStatusJSON(http.StatusForbidden, common.Forbidden.WithError(errors.New("forbidden: insufficient role")))

				return
			}
		}

		// todo: check disabled account
		// Save the account info in context
		c.Set(common.CurrentUserContextKey, account.User)
		c.Next()
	}
}
