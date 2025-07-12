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
		result, err2 := th.GetTodoDetailBz.GetTodoDetail(c, id)
		if err2 != nil {
			c.JSON(http.StatusInternalServerError, err2)
			return
		}
		result.MarkupId()
		c.JSON(http.StatusOK, common.SimpleResponse(result))
	}
}
