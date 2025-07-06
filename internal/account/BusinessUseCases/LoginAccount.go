package BusinessUseCases

import (
	"errors"
	"github.com/coderconquerer/go-login-app/internal/account/models"
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/coderconquerer/go-login-app/internal/components/tokenProviders"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FindAccountStorage interface {
	FindAccountByUsername(c *gin.Context, username string) (*models.Account, error)
}

type RegisterAccountBusiness interface {
	RegisterAccount(c *gin.Context, acc *models.Account) error
}

type LoginLogic struct {
	store         FindAccountStorage
	tokenProvider tokenProviders.TokenProvider
	expiry        int
}

func GetNewLoginLogic(store FindAccountStorage, tokenProvider tokenProviders.TokenProvider, expiry int) *LoginLogic {
	return &LoginLogic{
		store:         store,
		tokenProvider: tokenProvider,
		expiry:        expiry,
	}
}

func (bz *LoginLogic) Login(c *gin.Context, acc models.AccountLogin) (tokenProviders.Token, *common.AppError) {
	account, err := bz.store.FindAccountByUsername(c, acc.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewInvalidUsernameOrPassword("Incorrect username or password")
		} else {
			return nil, common.NewDatabaseError(err)
		}
	}

	isCorrect := common.ComparePasswordWithSalt(acc.Password, account.Salt, account.Password)
	if !isCorrect {
		return nil, common.NewInvalidUsernameOrPassword("Incorrect username or password")
	}
	payload := &common.TokenPayload{
		UserId: account.AccountID,
		Role:   account.Role,
	}

	accessToken, err := bz.tokenProvider.GenerateToken(payload, bz.expiry)
	if err != nil {
		return nil, common.NewInternalSeverErrorResponse(err, "error when creating access token", err.Error())
	}

	// todo: create refresh token
	return accessToken, nil

}
