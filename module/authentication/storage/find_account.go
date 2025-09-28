package storage

import (
	"context"
	"errors"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/authentication/entity"
	userEntity "github.com/coderconquerer/social-todo/module/user/entity"
	"go.opencensus.io/trace"
	"gorm.io/gorm"
)

func (db *MySQLConnection) FindAccountByUsername(c context.Context, username string) (*entity.Account, error) {
	_, span := trace.StartSpan(c, "user.storage.FindAccountByUsername")
	defer span.End()
	// filter deleted first
	var account entity.Account
	dbc := db.conn.Table(entity.Account{}.TableName())
	if err := dbc.Preload(userEntity.User{}.TableName()).Where("Username = ?", username).Take(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, common.DatabaseError.WithError(err)
	}

	return &account, nil
}

func (db *MySQLConnection) FindAccount(c context.Context, conditions map[string]interface{}) (*entity.Account, error) {
	_, span := trace.StartSpan(c, "user.storage.FindAccount")
	defer span.End()

	var account entity.Account
	dbc := db.conn.Table(entity.Account{}.TableName())

	if err := dbc.Preload(userEntity.User{}.TableName()).Where(conditions).Take(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}
