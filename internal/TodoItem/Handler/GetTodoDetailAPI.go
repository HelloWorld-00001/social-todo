package Handler

import (
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (th *TodoHandler) GetTodoDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewInvalidInputError(err))
			return
		}
		result, err := th.GetTodoDetailBz.GetTodoDetail(c, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		result.MarkupId()
		c.JSON(http.StatusOK, common.SimpleResponse(result))
	}
}
