package Handler

import (
	common2 "github.com/coderconquerer/social-todo/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (rh *ReactionHandler) GetListReactedUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

		todoId, err := common2.GetUidFromString(c.Param("todo_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, common2.NewInvalidInputError(err))
			return
		}

		pagination := common2.Pagination{}
		if err := c.ShouldBindQuery(&pagination); err != nil {
			c.JSON(http.StatusBadRequest, common2.NewInvalidInputError(err))
			return
		}
		pagination.Process()

		result, err2 := rh.ListReactedUserBz.GetListReactedUsers(c, int(todoId.LocalId()), &pagination)
		if err2 != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		for i := range result {
			result[i].CreateMarkupId()
		}
		c.JSON(http.StatusOK, common2.StandardResponseWithoutFilter(result, pagination))
	}
}
