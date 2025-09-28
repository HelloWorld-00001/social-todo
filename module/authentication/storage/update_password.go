package storage

import (
	"context"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/todo/entity"
)

func (db *MySQLConnection) UpdateTodoItem(c context.Context, filter *common.Filter, pagination *common.Pagination) ([]entity.Todo, error) {
	var todos []entity.Todo
	// filter deleted first
	dbc := db.conn.Table(entity.Todo{}.TableName())

	dbc = dbc.Where("deleted_date is null")
	if filter != nil {
		if stt := filter.Status; stt != "" {
			dbc.Where("status = ?", stt)
		}
	}

	if err := dbc.Select("id").Count(&pagination.Total).Error; err != nil {
		return nil, err
	}

	if err := dbc.Select("*").Order("id desc").
		Offset((pagination.Page - 1) * pagination.Limit).
		Limit(pagination.Limit).
		Find(&todos).Error; err != nil {
		return nil, err
	}

	return todos, nil
}
