package Handler

import (
	common2 "github.com/coderconquerer/go-login-app/common"
	"github.com/coderconquerer/go-login-app/module/account/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (ah *AccountHandler) RegisterAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var acc models.Account
		if err := c.ShouldBind(&acc); err != nil {
			c.JSON(http.StatusBadRequest, common2.NewInvalidInputError(err))
			return
		}
		err := ah.RegisterLogic.RegisterAccount(c, &acc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, common2.SimpleResponse(acc.AccountID))
	}
}
