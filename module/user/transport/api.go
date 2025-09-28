package userApi

import (
	"errors"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/user/business"
	"github.com/coderconquerer/social-todo/module/user/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserAPI interface {
	GetUserProfile() gin.HandlerFunc
}

type userAPI struct {
	business business.UserBusiness
}

func NewUserAPI(business business.UserBusiness) UserAPI {
	return &userAPI{
		business: business,
	}
}

func (uh *userAPI) GetUserProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		accInfo, ok := c.Get(common.CurrentUserContextKey)
		if !ok {
			common.RespondError(c,
				common.InternalServerError.
					WithError(errors.New("error when getting user information from given token")))
			return
		}

		username := accInfo.(*entity.User).Username

		userInfo, err := uh.business.GetUserProfileByUserName(c, username)
		if err != nil {
			common.RespondError(c, err)
			return
		}
		userInfo.CreateMarkupId()
		c.JSON(http.StatusOK, common.SimpleResponse(userInfo))
	}
}
