package Storage

import (
	"github.com/coderconquerer/go-login-app/internal/models"
	"github.com/gin-gonic/gin"
	"time"
)

func (db *MySQLConnection) DeleteTodoItem(c *gin.Context, id int) error {
	dbc := db.conn.Table(models.Todo{}.TableName())
	currentTime := time.Now()
	if err := dbc.Where("Id = ?", id).Updates(map[string]interface{}{
		"Deleted_Date": currentTime,
		"UpdateTime":   currentTime,
	}).Error; err != nil {
		return err
	}

	return nil
}
