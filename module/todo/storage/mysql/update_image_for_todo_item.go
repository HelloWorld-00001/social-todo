package mysql

import (
	"context"
	"database/sql/driver"
	"github.com/coderconquerer/social-todo/module/todo/entity"
)

func (db *MySQLConnection) UploadImageForTodo(c context.Context, id int, image driver.Value) error {
	// filter deleted first
	dbc := db.conn.Table(entity.Todo{}.TableName())

	dbc = dbc.Where("deleted_date is null")

	if err := dbc.Where("id = ?", id).Update("image", image).Error; err != nil {
		return err
	}

	return nil
}
