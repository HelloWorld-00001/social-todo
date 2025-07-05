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

func GetTodoDetail(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		store := Storage.GetNewMySQLConnection(db)
		business := BusinessUseCases.GetNewGetTodoDetailLogic(store)
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		result, err := business.GetTodoDetail(c, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, common.SimpleResponse(result))
	}
}
