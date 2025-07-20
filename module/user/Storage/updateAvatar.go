package Storage

import (
	"database/sql/driver"
	"github.com/coderconquerer/social-todo/module/user/models"
	"github.com/gin-gonic/gin"
)

func (db *MySQLConnection) UploadUserAvatar(c *gin.Context, id int, image driver.Value) error {
	// filter deleted first
	dbc := db.conn.Table(models.User{}.TableName())

	if err := dbc.Where("id = ?", id).Update("Avatar", image).Error; err != nil {
		return err
	}

	return nil
}
