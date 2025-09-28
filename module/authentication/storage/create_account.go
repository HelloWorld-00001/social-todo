package storage

import (
	"context"
	"github.com/coderconquerer/social-todo/module/authentication/entity"
)

func (db *MySQLConnection) CreateAccount(c context.Context, acc *entity.Account) error {
	// filter deleted first
	dbc := db.conn.Table(entity.Account{}.TableName())

	if err := dbc.Create(acc).Error; err != nil {
		return err
	}

	return nil
}
