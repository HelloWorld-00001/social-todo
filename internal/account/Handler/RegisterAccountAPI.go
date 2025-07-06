package Handler

import (
	"github.com/coderconquerer/go-login-app/internal/account/models"
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (ah *AccountHandler) RegisterAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var acc models.Account
		if err := c.ShouldBind(&acc); err != nil {
			c.JSON(http.StatusBadRequest, common.NewInvalidInputError(err))
			return
		}
		err := ah.RegisterLogic.RegisterAccount(c, &acc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, common.SimpleResponse(acc.AccountID))
	}
}
