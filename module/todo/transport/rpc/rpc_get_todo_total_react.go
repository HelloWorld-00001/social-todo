package rpc

import (
	"github.com/coderconquerer/social-todo/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (rh *RpcHandler) GetTodoTotalReact() gin.HandlerFunc {
	return func(c *gin.Context) {
		type TodoIds struct {
			Ids []int `json:"ids"`
		}
		var ids TodoIds
		if err := c.ShouldBind(&ids); err != nil {
			c.JSON(http.StatusBadRequest, common.NewInvalidInputError(err))
			return
		}

		reactionMap, err := rh.GetTodoItemTotalReactBz.GetTodoItemTotalReact(c, ids.Ids)

		if err != nil {
			c.JSON(http.StatusInternalServerError, common.NewInternalSeverErrorResponse(err, err.Error(), err.Error()))
			return
		}

		// Aggregate the reaction of to do
		c.JSON(http.StatusOK, common.SimpleResponse(reactionMap))
	}
}
