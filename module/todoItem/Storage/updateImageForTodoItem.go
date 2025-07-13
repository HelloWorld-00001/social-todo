package Storage

import (
	"database/sql/driver"
	"github.com/coderconquerer/go-login-app/module/todoItem/models"
	"github.com/gin-gonic/gin"
)

func (db *MySQLConnection) UploadImageForTodo(c *gin.Context, id int, image driver.Value) error {
	// filter deleted first
	dbc := db.conn.Table(models.Todo{}.TableName())

	dbc = dbc.Where("deleted_date is null")

	if err := dbc.Where("id = ?", id).Update("image", image).Error; err != nil {
		return err
	}

	return nil
}
