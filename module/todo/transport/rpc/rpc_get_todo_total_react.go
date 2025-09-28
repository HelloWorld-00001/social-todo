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
			common.RespondError(c, common.BadRequest.WithError(common.ErrInvalidInput).WithRootCause(err))
			return
		}

		reactionMap, errBz := rh.GetTodoItemTotalReactBz.GetTodoItemTotalReact(c, ids.Ids)

		if errBz != nil {
			common.RespondError(c, errBz)

			return
		}

		// Aggregate the reaction of to do
		c.JSON(http.StatusOK, common.SimpleResponse(reactionMap))
	}
}
