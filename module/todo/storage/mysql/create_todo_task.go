package mysql

import (
	"context"
	"github.com/coderconquerer/social-todo/module/todo/entity"
	"time"
)

func (db *MySQLConnection) CreateTodoItem(c context.Context, todo *entity.TodoCreation) error {
	// filter deleted first
	dbc := db.conn.Table(entity.Todo{}.TableName())
	currentTime := time.Now()
	todo.CreateTime = currentTime
	todo.UpdateTime = currentTime

	if err := dbc.Create(todo).Error; err != nil {
		return err
	}

	return nil
}
