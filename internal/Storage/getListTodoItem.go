package Storage

import (
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/coderconquerer/go-login-app/internal/models"
	"github.com/gin-gonic/gin"
)

func (db *MySQLConnection) GetTodoList(c *gin.Context, filter *common.Filter, pagination *common.Pagination) ([]models.Todo, error) {
	var todos []models.Todo
	// filter deleted first
	dbc := db.conn.Table(models.Todo{}.TableName())

	dbc = dbc.Where("Deleted_Date is null")
	if filter != nil {
		if stt := filter.Status; stt != "" {
			dbc.Where("Status = ?", stt)
		}
	}

	if err := dbc.Select("Id").Count(&pagination.Total).Error; err != nil {
		return nil, common.NewDatabaseError(err)
	}

	if err := dbc.Select("*").Order("Id desc").
		Offset((pagination.Page - 1) * pagination.Limit).
		Limit(pagination.Limit).
		Find(&todos).Error; err != nil {
		return nil, common.NewDatabaseError(err)
	}

	return todos, nil
}
