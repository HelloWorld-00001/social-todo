package Handler

import (
	common2 "github.com/coderconquerer/social-todo/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (th *TodoHandler) GetToDoList() gin.HandlerFunc {
	return func(c *gin.Context) {
		pagination := common2.Pagination{}
		if err := c.ShouldBindQuery(&pagination); err != nil {
			c.JSON(http.StatusBadRequest, common2.NewInvalidInputError(err))
			return
		}
		pagination.Process()

		result, err := th.GetTodoListBz.GetTodoList(c, nil, &pagination)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		for i := range result {
			result[i].CreateMarkupId()
		}
		c.JSON(http.StatusOK, common2.StandardResponseWithoutFilter(result, pagination))
	}
}
