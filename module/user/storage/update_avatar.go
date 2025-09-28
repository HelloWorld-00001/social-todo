package storage

import (
	"context"
	"database/sql/driver"
	"github.com/coderconquerer/social-todo/module/user/entity"
)

func (db *MySQLConnection) UploadUserAvatar(c context.Context, id int, image driver.Value) error {
	// filter deleted first
	dbc := db.conn.Table(entity.User{}.TableName())

	if err := dbc.Where("id = ?", id).Update("Avatar", image).Error; err != nil {
		return err
	}

	return nil
}
