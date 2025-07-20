package Handler

import (
	"errors"
	"github.com/coderconquerer/social-todo/common"
	model "github.com/coderconquerer/social-todo/module/user/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (rh *ReactionHandler) UnreactTodoItem() gin.HandlerFunc {
	return func(c *gin.Context) {

		uid, err := common.GetUidFromString(c.Query("todo_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewInvalidInputError(errors.New("Invalid todo id")))
			return
		}
		todoId := int(uid.LocalId())

		accInfo, ok := c.Get(common.CurrentUserContextKey)
		if !ok {
			c.JSON(http.StatusInternalServerError, common.NewInternalSeverErrorResponse(errors.New("cannot get user information"), "Error when getting user information", ""))
			return
		}

		err2 := rh.UnreactTodoItemBz.UnreactTodoItem(c, accInfo.(*model.User).Id, todoId)
		if err2 != nil {
			c.JSON(http.StatusInternalServerError, common.NewInternalSeverErrorResponse(err2, "Cannot unlike this item", ""))
			return
		}

		c.JSON(http.StatusOK, common.SimpleResponse(true))
	}
}
