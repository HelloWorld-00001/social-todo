package BusinessUseCases

import (
	"errors"
	common2 "github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/account/models"
	tokenProviders "github.com/coderconquerer/social-todo/plugin/tokenProviders"
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

func (bz *LoginLogic) Login(c *gin.Context, acc models.AccountLogin) (tokenProviders.Token, *common2.AppError) {
	account, err := bz.store.FindAccountByUsername(c, acc.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common2.NewInvalidUsernameOrPassword("Incorrect username or password")
		} else {
			return nil, common2.NewDatabaseError(err)
		}
	}
	if account.User == nil {
		return nil, common2.NewInternalSeverErrorResponse(errors.New("error when fetching User"), "", "")
	}

	isCorrect := common2.ComparePasswordWithSalt(acc.Password, account.Salt, account.Password)
	if !isCorrect {
		return nil, common2.NewInvalidUsernameOrPassword("Incorrect username or password")
	}
	payload := &common2.TokenPayload{
		AccountId: account.AccountID,
		UserId:    account.User.Id,
		Role:      account.Role,
	}

	accessToken, err := bz.tokenProvider.GenerateToken(payload, bz.expiry)
	if err != nil {
		return nil, common2.NewInternalSeverErrorResponse(err, "error when creating access token", err.Error())
	}

	// todo: create refresh token
	return accessToken, nil

}
