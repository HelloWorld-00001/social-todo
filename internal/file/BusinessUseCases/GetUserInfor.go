package BusinessUseCases

import (
	"errors"
	model "github.com/coderconquerer/go-login-app/internal/account/models"
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/coderconquerer/go-login-app/internal/user/models"
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

func (bz *FindUserLogic) GetUserProfile(c *gin.Context) (*models.User, *common.AppError) {
	accInfo, ok := c.Get(common.AccountContextKey)
	if !ok {
		return nil, common.NewInternalSeverErrorResponse(errors.New("cannot get user information"), "Error when getting user information", "")
	}
	cdt := map[string]interface{}{
		"Username": accInfo.(*model.Account).Username,
	}
	userInfo, err := bz.store.FindUser(c, cdt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewCannotGetEntity(models.User{}.TableName(), err)
		}
		return nil, common.NewDatabaseError(err)
	}

	return userInfo, nil
}
