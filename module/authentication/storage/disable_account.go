package storage

import (
	"context"
	"github.com/coderconquerer/social-todo/module/authentication/entity"
)

func (db *MySQLConnection) HandleDisableAccount(c context.Context, id int, isDisable bool) error {
	// filter deleted first
	dbc := db.conn.Table(entity.Account{}.TableName())

	if err := dbc.Where("id = ?", id).Update("IsDisable", isDisable).Error; err != nil {
		return err
	}

	return nil
}
