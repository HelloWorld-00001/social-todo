package Storage

import (
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/coderconquerer/go-login-app/internal/models"
	"github.com/gin-gonic/gin"
	"log"
)

func (db *MySQLConnection) GetTodoList(c *gin.Context, filter *common.Filter, pagination *common.Pagination) ([]models.Todo, error) {
	var todos []models.Todo
	// filter deleted first
	dbc := db.conn.Table(models.Todo{}.TableName())

	dbc = dbc.Where("DeletedDate is null")
	if filter != nil {
		if stt := filter.Status; stt != "" {
			dbc.Where("Status = ?", stt)
		}
	}

	if err := dbc.Select("Todo_Id").Count(&pagination.Total).Error; err != nil {
		return nil, err
	}
	log.Println("heheh")
	log.Println(pagination)

	if err := dbc.Select("*").Order("Todo_Id desc").
		Offset((pagination.Page - 1) * pagination.Limit).
		Limit(pagination.Limit).
		Find(&todos).Error; err != nil {
		return nil, err
	}

	return todos, nil
}
