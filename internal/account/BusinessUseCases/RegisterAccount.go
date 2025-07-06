package BusinessUseCases

import (
	"errors"
	"fmt"
	"github.com/coderconquerer/go-login-app/internal/account/models"
	"github.com/coderconquerer/go-login-app/internal/common"
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

func (bz *RegisterAccountLogic) RegisterAccount(c *gin.Context, acc *models.Account) *common.AppError {
	account, err := bz.store.FindAccountByUsername(c, acc.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return common.NewDatabaseError(err)
	}
	if account != nil {
		errMsg := fmt.Sprintf("Account with username %s already exists", acc.Username)
		return common.NewFullErrorResponse(errors.New(errMsg), errMsg, "", "Error_AccountConflict", http.StatusConflict)
	}
	salt, err := common.GenerateSalt(64)
	if err != nil {
		errMsg := fmt.Sprintf("Salt generation failed: %s", err.Error())
		return common.NewInternalSeverErrorResponse(err, "Error when generating salt", errMsg)
	}
	pw, err := common.HashPasswordWithSalt(acc.Password, salt)
	if err != nil {
		errMsg := fmt.Sprintf("Error when hashing password: %s", err.Error())
		return common.NewInternalSeverErrorResponse(err, "Error when hashing password", errMsg)
	}
	acc.Password = pw
	acc.Salt = salt
	acc.Role = common.UserRole.ToString()
	if err = bz.store.CreateAccount(c, acc); err != nil {
		return common.NewCannotCreateEntity(models.Account{}.TableName(), err)
	}
	return nil
}
