package BusinessUseCases

import (
	"errors"
	"fmt"
	common2 "github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/account/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type CreateAccountStorage interface {
	CreateAccount(c *gin.Context, acc *models.Account) error
	FindAccountByUsername(c *gin.Context, username string) (*models.Account, error)
}

type RegisterAccountLogic struct {
	store CreateAccountStorage
}

func GetNewRegisterAccountLogic(store CreateAccountStorage) *RegisterAccountLogic {
	return &RegisterAccountLogic{store: store}
}

func (bz *RegisterAccountLogic) RegisterAccount(c *gin.Context, acc *models.Account) *common2.AppError {
	account, err := bz.store.FindAccountByUsername(c, acc.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return common2.NewDatabaseError(err)
	}
	if account != nil {
		errMsg := fmt.Sprintf("Account with username %s already exists", acc.Username)
		return common2.NewFullErrorResponse(errors.New(errMsg), errMsg, "", "Error_AccountConflict", http.StatusConflict)
	}
	salt, err := common2.GenerateSalt(64)
	if err != nil {
		errMsg := fmt.Sprintf("Salt generation failed: %s", err.Error())
		return common2.NewInternalSeverErrorResponse(err, "Error when generating salt", errMsg)
	}
	pw, err := common2.HashPasswordWithSalt(acc.Password, salt)
	if err != nil {
		errMsg := fmt.Sprintf("Error when hashing password: %s", err.Error())
		return common2.NewInternalSeverErrorResponse(err, "Error when hashing password", errMsg)
	}
	acc.Password = pw
	acc.Salt = salt
	acc.Role = common2.UserRole.ToString()
	if err = bz.store.CreateAccount(c, acc); err != nil {
		return common2.NewCannotCreateEntity(models.Account{}.TableName(), err)
	}
	return nil
}
