package Handler

import (
	"github.com/coderconquerer/social-todo/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (ah *UserHandler) DisableAccount() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.JSON(http.StatusOK, common.SimpleResponse(true))
	}
}
