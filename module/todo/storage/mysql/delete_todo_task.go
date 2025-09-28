package mysql

import (
	"context"
	"github.com/coderconquerer/social-todo/module/todo/entity"
	"time"
)

func (db *MySQLConnection) DeleteTodoItem(c context.Context, id int) error {
	dbc := db.conn.Table(entity.Todo{}.TableName())
	currentTime := time.Now()
	if err := dbc.Where("Id = ?", id).Updates(map[string]interface{}{
		"Deleted_Date": currentTime,
		"UpdateTime":   currentTime,
	}).Error; err != nil {
		return err
	}

	return nil
}
