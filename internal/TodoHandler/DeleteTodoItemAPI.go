package TodoHandler

import (
	"github.com/coderconquerer/go-login-app/internal/BusinessUseCases"
	"github.com/coderconquerer/go-login-app/internal/Storage"
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		if err := business.DeleteTodoItem(c, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, common.SimpleResponse(true))
	}
}
