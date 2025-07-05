package TodoHandler

import (
	"github.com/coderconquerer/go-login-app/internal/BusinessUseCases"
	"github.com/coderconquerer/go-login-app/internal/Storage"
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/coderconquerer/go-login-app/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func CreateTodoItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		store := Storage.GetNewMySQLConnection(db)
		business := BusinessUseCases.GetNewCreateTodoLogic(store)
		var todo models.TodoCreation
		if err := c.ShouldBind(&todo); err != nil {
			c.JSON(http.StatusBadRequest, common.NewInvalidInputError(err))
			return
		}
		err := business.CreateTodoItem(c, &todo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, common.SimpleResponse(todo.TodoID))
	}
}
