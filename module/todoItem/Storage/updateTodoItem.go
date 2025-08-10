package Storage

import (
	"context"
	common2 "github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/todoItem/models"
	"gorm.io/gorm"
)

func (db *MySQLConnection) UpdateTodoItem(c context.Context, filter *common2.Filter, pagination *common2.Pagination) ([]models.Todo, error) {
	var todos []models.Todo
	// filter deleted first
	dbc := db.conn.Table(models.Todo{}.TableName())

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

func (db *MySQLConnection) IncreaseTotalReactionCount(c context.Context, todoId int) error {
	// filter deleted first
	dbc := db.conn.Table(models.Todo{}.TableName())

	if err := dbc.Where("Id = ?", todoId).
		UpdateColumn("TotalReact", gorm.Expr("TotalReact + ?", 1)).
		Error; err != nil {
		return err
	}

	return nil
}

func (db *MySQLConnection) DecreaseTotalReactionCount(c context.Context, todoId int) error {
	// filter deleted first
	dbc := db.conn.Table(models.Todo{}.TableName())

	if err := dbc.Where("Id = ?", todoId).
		UpdateColumn("TotalReact", gorm.Expr("TotalReact - ?", 1)).
		Error; err != nil {
		return err
	}

	return nil
}
