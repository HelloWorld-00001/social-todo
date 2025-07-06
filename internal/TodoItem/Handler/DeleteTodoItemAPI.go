package Handler

import (
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (th *TodoHandler) DeleteTodoItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewInvalidInputError(err))
			return
		}

		if err := th.DeleteTodoBz.DeleteTodoItem(c, id); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, common.SimpleResponse(true))
	}
}
