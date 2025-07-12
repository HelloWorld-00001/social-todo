package Handler

import (
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (th *TodoHandler) GetToDoList() gin.HandlerFunc {
	return func(c *gin.Context) {
		pagination := common.Pagination{}
		if err := c.ShouldBindQuery(&pagination); err != nil {
			c.JSON(http.StatusBadRequest, common.NewInvalidInputError(err))
			return
		}
		log.Println(pagination)
		pagination.Process()

		result, err := th.GetTodoListBz.GetTodoList(c, nil, &pagination)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		for i := range result {
			result[i].MarkupId() // ✅ modifies original struct
		}
		c.JSON(http.StatusOK, common.StandardResponseWithoutFilter(result, pagination))
	}
}
