package Handler

import (
	"github.com/coderconquerer/go-login-app/internal/account/models"
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (ah *UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var acc models.AccountLogin
		if err := c.ShouldBind(&acc); err != nil {
			c.JSON(http.StatusBadRequest, common.NewInvalidInputError(err))
			return
		}
		//result, err := ah.LoginLogic.Login(c, acc)
		//if err != nil {
		//	c.JSON(http.StatusInternalServerError, err)
		//	return
		////}
		//if result == nil {
		//	errMsg := "username or password is incorrect"
		//	c.JSON(http.StatusBadRequest, common.NewFullErrorResponse(errors.New(errMsg), errMsg, "", "Error_IncorrectLogin", http.StatusBadRequest))
		//	return
		//}

		c.JSON(http.StatusOK, common.SimpleResponse(true))
	}
}
