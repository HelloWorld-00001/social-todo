package Handler

import (
	"github.com/coderconquerer/go-login-app/internal/TodoItem/BusinessUseCases"
	"github.com/coderconquerer/go-login-app/internal/TodoItem/Storage"
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func DeleteTodoItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		store := Storage.GetNewMySQLConnection(db)
		business := BusinessUseCases.GetNewDeleteTodoItemLogic(store)
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewInvalidInputError(err))
			return
		}

		if err := business.DeleteTodoItem(c, id); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, common.SimpleResponse(true))
	}
}
