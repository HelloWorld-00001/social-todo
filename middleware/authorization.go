package middleware

import (
	"context"
	"errors"
	common2 "github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/account/models"
	tokenProviders "github.com/coderconquerer/social-todo/plugin/tokenProviders"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type AuthorizationStore interface {
	FindAccount(c context.Context, conditions map[string]interface{}) (*models.Account, error)
}

func ErrorWrongAuthHeader(err error) *common2.AppError {
	return common2.NewBadRequestResponseWithError(err, "Wrong authorization header", err.Error())
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

		payload, err := provider.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common2.NewUnauthorizedErrorCustom(err, "Invalid or expired token"))
			return
		}

		condition := map[string]interface{}{
			"Id": payload.GetAccountId(),
		}

		ctx := c.Request.Context()
		account, err := store.FindAccount(ctx, condition)
		if err != nil || account == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common2.NewUnauthorizedErrorCustom(err, "Account not found"))
			return
		}

		if account.User == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common2.NewUnauthorizedErrorCustom(err, "Cannot get user information"))
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
				c.AbortWithStatusJSON(http.StatusForbidden, common2.NewForbiddenError("Forbidden: insufficient role"))
				return
			}
		}

		// todo: check disabled account
		// Save the account info in context
		c.Set(common2.CurrentUserContextKey, account.User)
		c.Next()
	}
}
