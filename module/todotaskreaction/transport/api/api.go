package api

import (
	"errors"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/todotaskreaction/business"
	"github.com/coderconquerer/social-todo/module/todotaskreaction/entity"
	userEntity "github.com/coderconquerer/social-todo/module/user/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ReactionTodoAPI interface {
	ReactItem() gin.HandlerFunc
	UnreactTodoItem() gin.HandlerFunc
	GetListReactedUsers() gin.HandlerFunc
}

type reactionTodoAPI struct {
	business business.ReactionBusiness
}

func NewReactionTodoAPI(business business.ReactionBusiness) ReactionTodoAPI {
	return &reactionTodoAPI{business}
}

func (rh *reactionTodoAPI) ReactItem() gin.HandlerFunc {
	return func(c *gin.Context) {

		var input entity.ReactionInput
		if err := c.ShouldBindQuery(&input); err != nil {
			common.RespondError(c,
				common.BadRequest.
					WithError(common.ErrInvalidInput).
					WithRootCause(err),
			)
			return
		}

		ids, errParsingId := common.GetUidFromString(input.TodoId)
		if errParsingId != nil {
			common.RespondError(c,
				common.BadRequest.
					WithError(common.ErrInvalidInput).
					WithRootCause(errParsingId),
			)
			return
		}
		id := int(ids.LocalId())

		reactEnum, errReaction := common.GetReactionFromString(input.React)
		if errReaction != nil {
			common.RespondError(c,
				common.BadRequest.
					WithError(common.ErrInvalidInput).
					WithRootCause(errReaction),
			)
			return
		}

		accInfo, ok := c.Get(common.CurrentUserContextKey)
		if !ok {
			common.RespondError(c,
				common.InternalServerError.
					WithError(errors.New("error when getting user information from given token")))
			return
		}

		errBz := rh.business.ReactTodoItem(c, entity.Reaction{
			UserId:    accInfo.(*userEntity.User).Id,
			TodoId:    id,
			CreatedAt: time.Now(),
			React:     reactEnum,
		})
		if errBz != nil {
			common.RespondError(c, errBz)
			return
		}

		c.JSON(http.StatusOK, common.SimpleResponse(true))
	}
}

func (rh *reactionTodoAPI) UnreactTodoItem() gin.HandlerFunc {
	return func(c *gin.Context) {

		uid, errParsingId := common.GetUidFromString(c.Query("todo_id"))
		if errParsingId != nil {
			common.RespondError(c,
				common.BadRequest.
					WithError(common.ErrInvalidInput).
					WithRootCause(errParsingId),
			)
			return
		}
		todoId := int(uid.LocalId())

		accInfo, ok := c.Get(common.CurrentUserContextKey)
		if !ok {
			common.RespondError(c,
				common.InternalServerError.
					WithError(errors.New("error when getting user information from given token")))
			return
		}

		err := rh.business.UnreactTodoItem(c, accInfo.(*userEntity.User).Id, todoId)
		if err != nil {
			common.RespondError(c, err)
			return
		}

		c.JSON(http.StatusOK, common.SimpleResponse(true))
	}
}

func (rh *reactionTodoAPI) GetListReactedUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

		todoId, errParsingId := common.GetUidFromString(c.Param("todo_id"))
		if errParsingId != nil {
			common.RespondError(c,
				common.BadRequest.
					WithError(common.ErrInvalidInput).
					WithRootCause(errParsingId),
			)
			return
		}

		pagination := common.Pagination{}
		if errParsingPagination := c.ShouldBindQuery(&pagination); errParsingPagination != nil {
			common.RespondError(c,
				common.BadRequest.
					WithError(common.ErrInvalidInput).
					WithRootCause(errParsingPagination),
			)
			return
		}
		pagination.Process()

		result, err := rh.business.GetListReactedUsers(c, int(todoId.LocalId()), &pagination)
		if err != nil {
			common.RespondError(c, err)
			return
		}

		for i := range result {
			result[i].CreateMarkupId()
		}
		c.JSON(http.StatusOK, common.StandardResponseWithoutFilter(result, pagination))
	}
}
