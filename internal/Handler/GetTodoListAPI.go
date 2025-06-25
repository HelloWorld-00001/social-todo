package Handler

import (
	"github.com/coderconquerer/go-login-app/internal/BusinessUseCases"
	"github.com/coderconquerer/go-login-app/internal/Storage"
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func GetToDoList(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		store := Storage.GetNewMySQLConnection(db)
		business := BusinessUseCases.GetNewGetTodoListLogic(store)
		pagination := common.Pagination{}
		if err := c.ShouldBindQuery(&pagination); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Println(pagination)
		pagination.Process()

		result, err := business.GetTodoList(c, nil, &pagination)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, common.StandardResponseWithoutFilter(result, pagination))
	}
}
