package Storage

import (
	"github.com/coderconquerer/go-login-app/module/user/models"
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
