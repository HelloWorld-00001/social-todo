package Handler

import (
	"github.com/coderconquerer/go-login-app/internal/TodoItem/BusinessUseCases"
	"github.com/coderconquerer/go-login-app/internal/TodoItem/Storage"
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
			c.JSON(http.StatusBadRequest, common.NewInvalidInputError(err))
			return
		}
		log.Println(pagination)
		pagination.Process()

		result, err := business.GetTodoList(c, nil, &pagination)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, common.StandardResponseWithoutFilter(result, pagination))
	}
}
