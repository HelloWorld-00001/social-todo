package Handler

import (
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/coderconquerer/go-login-app/internal/todoItem/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (th *TodoHandler) CreateTodoItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var todo models.TodoCreation
		if err := c.ShouldBind(&todo); err != nil {
			c.JSON(http.StatusBadRequest, common.NewInvalidInputError(err))
			return
		}
		err := th.CreateTodoBz.CreateTodoItem(c, &todo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, common.SimpleResponse(todo.TodoID))
	}
}
