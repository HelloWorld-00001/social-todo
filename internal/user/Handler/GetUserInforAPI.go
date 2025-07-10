package Handler

import (
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (uh *UserHandler) GetUserProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, err := uh.GetUserBz.GetUserProfile(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		userInfo.MarkupId()
		c.JSON(http.StatusOK, common.SimpleResponse(userInfo))
	}
}
