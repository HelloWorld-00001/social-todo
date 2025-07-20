package Storage

import (
	"github.com/coderconquerer/social-todo/module/user/models"
	"github.com/gin-gonic/gin"
)

func (db *MySQLConnection) CreateUser(c *gin.Context, acc *models.User) error {
	// filter deleted first
	dbc := db.conn.Table(models.User{}.TableName())

	if err := dbc.Create(acc).Error; err != nil {
		return err
	}

	return nil
}
