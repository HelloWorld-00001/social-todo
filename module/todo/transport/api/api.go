package api

import (
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/todo/business"
	"github.com/coderconquerer/social-todo/module/todo/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TodoAPI interface {
	CreateTodoItem() gin.HandlerFunc
	DeleteTodoItem() gin.HandlerFunc
	GetTodoDetail() gin.HandlerFunc
	GetToDoList() gin.HandlerFunc
}
type todoAPI struct {
	business business.TodoBusiness
}

func NewTodoAPI(business business.TodoBusiness) TodoAPI {
	return &todoAPI{
		business,
	}
}

func (th *todoAPI) CreateTodoItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var todo entity.TodoCreation
		if err := c.ShouldBind(&todo); err != nil {
			common.RespondError(c,
				common.BadRequest.
					WithError(entity.ErrInvalidTodoCreation).
					WithRootCause(err),
			)
			return
		}
		err := th.business.CreateTodoItem(c, &todo)
		if err != nil {
			common.RespondError(c, err)
			return
		}

		c.JSON(http.StatusOK, common.SimpleResponse(todo.TodoID))
	}
}

func (th *todoAPI) DeleteTodoItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			common.RespondError(c,
				common.BadRequest.
					WithError(common.ErrInvalidInputID).
					WithRootCause(err),
			)
			return
		}

		if err := th.business.DeleteTodoItem(c, id); err != nil {
			common.RespondError(c, err)
			return
		}

		c.JSON(http.StatusOK, common.SimpleResponse(true))
	}
}

func (th *todoAPI) GetTodoDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, inputErr := common.GetUidFromString(c.Param("id"))
		if inputErr != nil {
			common.RespondError(c,
				common.BadRequest.
					WithError(common.ErrInvalidInputID).
					WithRootCause(inputErr),
			)
			return
		}
		result, err := th.business.GetTodoDetail(c, int(id.LocalId()))
		if err != nil {
			common.RespondError(c, err)
			return
		}
		result.CreateMarkupId()
		c.JSON(http.StatusOK, common.SimpleResponse(result))
	}
}

func (th *todoAPI) GetToDoList() gin.HandlerFunc {
	return func(c *gin.Context) {
		pagination := common.Pagination{}
		if err := c.ShouldBindQuery(&pagination); err != nil {
			common.RespondError(c,
				common.BadRequest.
					WithError(common.ErrInvalidInputPagination).
					WithRootCause(err),
			)
			return
		}
		pagination.Process()

		result, err := th.business.GetTodoList(c.Request.Context(), nil, &pagination)
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
