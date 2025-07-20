package Handler

import (
	common2 "github.com/coderconquerer/social-todo/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (th *TodoHandler) DeleteTodoItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, common2.NewInvalidInputError(err))
			return
		}

		if err := th.DeleteTodoBz.DeleteTodoItem(c, id); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, common2.SimpleResponse(true))
	}
}
