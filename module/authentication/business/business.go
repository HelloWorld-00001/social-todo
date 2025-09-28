package business

import (
	"context"
	"errors"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/authentication/entity"
	tokenProviders "github.com/coderconquerer/social-todo/plugin/tokenProviders"
	"gorm.io/gorm"
)

//
// Unified Account Storage Interface
//

type AuthenticationStorage interface {
	CreateAccount(c context.Context, acc *entity.Account) error
	FindAccount(c context.Context, conditions map[string]interface{}) (*entity.Account, error)
	FindAccountByUsername(c context.Context, username string) (*entity.Account, error)
	HandleDisableAccount(c context.Context, id int, isDisable bool) error
}

// AuthenticationBusiness defines the contract for authentication logic
type AuthenticationBusiness interface {
	Login(c context.Context, acc entity.AccountLogin) (tokenProviders.Token, error)
	RegisterAccount(c context.Context, acc *entity.AccountRegister) error
	DisableAccount(c context.Context, id int, isDisable bool) error
}

type authenticationBusiness struct {
	store         AuthenticationStorage
	tokenProvider tokenProviders.TokenProvider
	expiry        int
}

func NewAuthenticationBusiness(store AuthenticationStorage, tokenProvider tokenProviders.TokenProvider, expiry int) *authenticationBusiness {
	return &authenticationBusiness{store: store, tokenProvider: tokenProvider, expiry: expiry}
}

func (bz *authenticationBusiness) Login(c context.Context, acc entity.AccountLogin) (tokenProviders.Token, error) {
	account, err := bz.store.FindAccountByUsername(c, acc.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.BadRequest.WithError(entity.ErrCannotFindAccount).WithRootCause(err)
		}
		return nil, common.InternalServerError.WithError(common.ErrUnhandleError).WithRootCause(err)
	}

	if account.User == nil {
		return nil, common.InternalServerError.WithError(entity.ErrCannotFindAccount)
	}

	if !common.ComparePasswordWithSalt(acc.Password, account.Salt, account.Password) {
		return nil, common.BadRequest.WithError(entity.ErrInvalidLoginCredential).WithRootCause(err)
	}

	payload := &common.TokenPayload{
		AccountId: account.Id,
		UserId:    account.User.Id,
		Role:      account.Role,
	}

	accessToken, err := bz.tokenProvider.GenerateToken(payload, bz.expiry)
	if err != nil {
		return nil, common.InternalServerError.WithError(errors.New("cannot generate token")).WithRootCause(err)
	}

	// TODO: implement refresh token
	return accessToken, nil
}

func (bz *authenticationBusiness) RegisterAccount(c context.Context, acc *entity.AccountRegister) error {
	account, err := bz.store.FindAccountByUsername(c, acc.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return common.InternalServerError.WithError(common.ErrUnhandleError).WithRootCause(err)
	}
	if account != nil {
		return common.Conflict.WithError(entity.ErrAccountConflict)
	}

	salt, errSalt := common.GenerateSalt(64)
	if errSalt != nil {
		return common.InternalServerError.WithError(errors.New("error when generating salt for password")).WithRootCause(errSalt)
	}

	pw, errHash := common.HashPasswordWithSalt(acc.Password, salt)
	if errHash != nil {
		return common.InternalServerError.WithError(errors.New("error when hashing password")).WithRootCause(errHash)
	}

	newAccount := &entity.Account{
		Username: acc.Username,
		Password: pw,
		Salt:     salt,
		Role:     common.UserRole.ToString(),
	}

	if errBz := bz.store.CreateAccount(c, newAccount); errBz != nil {
		return common.InternalServerError.WithError(entity.ErrCannotCreateAccount).WithRootCause(errBz)
	}
	return nil
}

func (bz *authenticationBusiness) DisableAccount(c context.Context, id int, isDisable bool) error {
	cdt := map[string]interface{}{"id": id}

	_, err := bz.store.FindAccount(c, cdt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.BadRequest.WithError(entity.ErrCannotFindAccount).WithRootCause(err)
		}
		return common.InternalServerError.WithError(entity.ErrCannotFindAccount)
	}

	if errCreate := bz.store.HandleDisableAccount(c, id, isDisable); errCreate != nil {
		return common.InternalServerError.WithError(entity.ErrCannotCreateAccount).WithRootCause(errCreate)
	}
	return nil
}
