package transport

import (
	"errors"
	"github.com/coderconquerer/social-todo/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"github.com/coderconquerer/social-todo/module/authenticationrpc/business"
	"github.com/coderconquerer/social-todo/module/authenticationrpc/entity"
)

type AuthenticationRpcAPI interface {
	Login() gin.HandlerFunc
	DisableAccount() gin.HandlerFunc
	RegisterAccount() gin.HandlerFunc
}

type authenticationAPI struct {
	business business.AuthenticationBusinessGrpc
}

func NewAuthenticationAPI(logic business.AuthenticationBusinessGrpc) AuthenticationRpcAPI {
	return &authenticationAPI{logic}
}

func (ah *authenticationAPI) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var acc entity.AccountLogin
		if err := c.ShouldBind(&acc); err != nil {
			common.RespondError(c,
				common.BadRequest.
					WithError(common.ErrInvalidInput).
					WithRootCause(err),
			)
			return
		}
		result, err := ah.business.Login(c, acc)
		if err != nil {
			common.RespondError(c, err)
			return
		}
		if result == "" {
			common.RespondError(c,
				common.BadRequest.
					WithError(entity.ErrInvalidLoginCredential).
					WithRootCause(err),
			)
			return
		}

		c.JSON(http.StatusOK, common.SimpleResponse(result))
	}
}

func (ah *authenticationAPI) RegisterAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var acc entity.AccountRegister
		if err := c.ShouldBind(&acc); err != nil {
			common.RespondError(c,
				common.BadRequest.
					WithError(common.ErrInvalidInput).
					WithRootCause(err),
			)
			return
		}
		err := ah.business.RegisterAccount(c, &acc)
		if err != nil {
			common.RespondError(c, err)
			return
		}

		c.JSON(http.StatusOK, common.SimpleResponse(true))
	}
}

func (ah *authenticationAPI) DisableAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract query parameters
		id := c.Query("id")           // returns string (empty if missing)
		disable := c.Query("disable") // returns string

		// Optional: convert to int
		idInt, err := strconv.Atoi(id)
		if err != nil {
			common.RespondError(c,
				common.BadRequest.
					WithError(common.ErrInvalidInput).
					WithRootCause(err),
			)
			return
		}
		disableInt, err := strconv.Atoi(disable)
		if err != nil {
			common.RespondError(c,
				common.BadRequest.
					WithError(common.ErrInvalidInput).
					WithRootCause(err),
			)
			return
		}
		if disableInt != 0 && disableInt != 1 {
			common.RespondError(c,
				common.BadRequest.
					WithError(errors.New("disable must be 0 or 1")).
					WithRootCause(err),
			)
			return
		}

		if errDisable := ah.business.DisableAccount(c, idInt, disableInt); errDisable != nil {
			common.RespondError(c, errDisable)
			return
		}

		c.JSON(http.StatusOK, common.SimpleResponse(true))
	}
}
