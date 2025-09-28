package business

import (
	"context"
	"errors"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/user/entity"
	"gorm.io/gorm"
)

type UserStorage interface {
	FindUserById(c context.Context, id int) (*entity.User, error)
	FindUser(c context.Context, conditions map[string]interface{}) (*entity.User, error)
}

type UserBusiness interface {
	GetUserProfileByUserName(c context.Context, username string) (*entity.User, error)
}
type userBusiness struct {
	store UserStorage
}

func NewUserBusiness(store UserStorage) UserBusiness {
	return &userBusiness{store: store}
}

func (bz *userBusiness) GetUserProfileByUserName(c context.Context, username string) (*entity.User, error) {
	cdt := map[string]interface{}{
		"Username": username,
	}
	userInfo, err := bz.store.FindUser(c, cdt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NotFound.WithError(entity.ErrCannotFindUser).WithRootCause(err)
		}
		return nil, common.InternalServerError.WithError(common.ErrUnhandleError).WithRootCause(err)
	}

	return userInfo, nil
}
