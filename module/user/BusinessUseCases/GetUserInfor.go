package BusinessUseCases

import (
	"errors"
	common2 "github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/user/models"
	model "github.com/coderconquerer/social-todo/module/user/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FindUserStorage interface {
	FindUserById(c *gin.Context, id int) (*models.User, error)
	FindUser(c *gin.Context, conditions map[string]interface{}) (*models.User, error)
}

type FindUserLogic struct {
	store FindUserStorage
}

func GetNewFindUserLogic(store FindUserStorage) *FindUserLogic {
	return &FindUserLogic{store: store}
}

func (bz *FindUserLogic) GetUserProfile(c *gin.Context) (*models.User, *common2.AppError) {
	accInfo, ok := c.Get(common2.CurrentUserContextKey)
	if !ok {
		return nil, common2.NewInternalSeverErrorResponse(errors.New("cannot get user information"), "Error when getting user information", "")
	}
	cdt := map[string]interface{}{
		"Username": accInfo.(*model.User).Username,
	}
	userInfo, err := bz.store.FindUser(c, cdt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common2.NewCannotGetEntity(models.User{}.TableName(), err)
		}
		return nil, common2.NewDatabaseError(err)
	}

	return userInfo, nil
}
