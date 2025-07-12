package Storage

import (
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/coderconquerer/go-login-app/internal/todoItem/models"
	"github.com/gin-gonic/gin"
)

func (db *MySQLConnection) UpdateTodoItem(c *gin.Context, filter *common.Filter, pagination *common.Pagination) ([]models.Todo, error) {
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
