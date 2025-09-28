package mysql

import (
	"context"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/todo/entity"
	"go.opencensus.io/trace"
)

func (db *MySQLConnection) GetTodoList(c context.Context, filter *common.Filter, pagination *common.Pagination) ([]entity.Todo, error) {

	_, span := trace.StartSpan(c, "todo.storage.GetTodoList")
	defer span.End()

	var todos []entity.Todo
	// filter deleted first
	dbc := db.conn.Table(entity.Todo{}.TableName())

	dbc = dbc.Where("Deleted_Date is null")
	if filter != nil {
		if stt := filter.Status; stt != "" {
			dbc.Where("Status = ?", stt)
		}
	}

	if err := dbc.Select("Id").Count(&pagination.Total).Error; err != nil {
		return nil, err
	}

	if pagination.Cursor != "" {
		cursor, err := common.GetUidFromString(pagination.Cursor)
		if err != nil {
			return nil, err
		}

		dbc.Where("Id < ?", cursor.LocalId())
	} else {
		dbc.Offset((pagination.Page - 1) * pagination.Limit)
	}

	if err := dbc.Select("*").Order("Id desc").
		Limit(pagination.Limit).
		Find(&todos).Error; err != nil {
		return nil, err
	}

	size := len(todos)
	if size == pagination.Limit {
		todos[size-1].CreateMarkupId()
		pagination.NextCursor = todos[size-1].MarkupId.String()
	}
	return todos, nil
}
