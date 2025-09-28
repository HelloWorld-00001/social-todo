package mysql

import (
	"context"
	"github.com/coderconquerer/social-todo/module/todo/entity"
)

func (db *MySQLConnection) GetTodoItemDetailById(c context.Context, id int) (*entity.Todo, error) {
	var todo entity.Todo
	// filter deleted first
	dbc := db.conn.Table(entity.Todo{}.TableName())

	// Todo: if account's role = admin -> return item despite of being deleted or not,
	// otherwise, return not found if record is deleted
	dbc = dbc.Where("Deleted_Date is null")

	if err := dbc.Where("Id = ?", id).Take(&todo).Error; err != nil {
		return nil, err
	}

	return &todo, nil
}
