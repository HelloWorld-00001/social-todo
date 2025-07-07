package BusinessUseCases

import (
	"errors"
	"github.com/coderconquerer/go-login-app/internal/account/models"
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DisableAccountStorage interface {
	HandleDisableAccount(c *gin.Context, id int, isDisable bool) error
	FindAccount(c *gin.Context, conditions map[string]interface{}) (*models.Account, error)
}

type DisableAccountLogic struct {
	store DisableAccountStorage
}

func GetNewDisableAccountLogic(store DisableAccountStorage) *DisableAccountLogic {
	return &DisableAccountLogic{store: store}
}

func (bz *DisableAccountLogic) DisableAccount(c *gin.Context, id int, isDisable bool) *common.AppError {
	cdt := map[string]interface{}{
		"id": id,
	}
	_, err := bz.store.FindAccount(c, cdt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NewNotFoundErrorResponse(err)
		} else {
			return common.NewDatabaseError(err)
		}
	}
	err = bz.store.HandleDisableAccount(c, id, isDisable)
	if err != nil {
		return common.NewCannotUpdateEntity(models.Account{}.TableName(), err)
	}
	return nil
}
