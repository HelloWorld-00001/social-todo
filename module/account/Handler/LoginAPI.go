package Handler

import (
	"errors"
	common2 "github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/account/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (ah *AccountHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var acc models.AccountLogin
		if err := c.ShouldBind(&acc); err != nil {
			c.JSON(http.StatusBadRequest, common2.NewInvalidInputError(err))
			return
		}
		result, err := ah.LoginLogic.Login(c, acc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		if result == nil {
			errMsg := "username or password is incorrect"
			c.JSON(http.StatusBadRequest, common2.NewFullErrorResponse(errors.New(errMsg), errMsg, "", "Error_IncorrectLogin", http.StatusBadRequest))
			return
		}

		c.JSON(http.StatusOK, common2.SimpleResponse(result))
	}
}
