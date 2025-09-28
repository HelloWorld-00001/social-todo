package storage

import (
	"context"
	"errors"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/user/entity"
	"gorm.io/gorm"
)

func (db *MySQLConnection) FindUserById(c context.Context, id int) (*entity.User, error) {
	// filter deleted first
	var user entity.User
	dbc := db.conn.Table(entity.User{}.TableName())
	if err := dbc.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, common.DatabaseError.WithError(err)
	}

	return &user, nil
}

func (db *MySQLConnection) FindUser(c context.Context, conditions map[string]interface{}) (*entity.User, error) {
	var user entity.User
	dbc := db.conn.Table(entity.User{}.TableName())

	if err := dbc.Where(conditions).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
