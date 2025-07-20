package Handler

import (
	"errors"
	"github.com/coderconquerer/social-todo/common"
	model "github.com/coderconquerer/social-todo/module/user/models"
	"github.com/coderconquerer/social-todo/module/userReactItem/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (rh *ReactionHandler) ReactItem() gin.HandlerFunc {
	return func(c *gin.Context) {

		var input models.ReactionInput
		if err := c.ShouldBindQuery(&input); err != nil {
			c.JSON(http.StatusBadRequest, common.NewInvalidInputError(err))
			return
		}

		ids, err := common.GetUidFromString(input.TodoId)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewInvalidInputError(err))
			return
		}
		id := int(ids.LocalId())

		reactEnum, err := common.GetReactionFromString(input.React)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewInvalidInputError(errors.New("Invalid reaction type")))
			return
		}

		accInfo, ok := c.Get(common.CurrentUserContextKey)
		if !ok {
			c.JSON(http.StatusInternalServerError, common.NewInternalSeverErrorResponse(errors.New("cannot get user information"), "Error when getting user information", ""))
			return
		}

		err2 := rh.ReactTodoItemBz.ReactTodoItem(c, models.Reaction{
			UserId:    accInfo.(*model.User).Id,
			TodoId:    id,
			CreatedAt: time.Now(),
			React:     reactEnum,
		})
		if err2 != nil {
			c.JSON(http.StatusInternalServerError, common.NewInternalSeverErrorResponse(err2, "", ""))
			return
		}

		c.JSON(http.StatusOK, common.SimpleResponse(true))
	}
}
