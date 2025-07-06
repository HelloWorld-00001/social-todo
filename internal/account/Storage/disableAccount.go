package Storage

import (
	"github.com/coderconquerer/go-login-app/internal/account/models"
	"github.com/gin-gonic/gin"
)

func (db *MySQLConnection) HandleDisableAccount(c *gin.Context, id int, isDisable bool) error {
	// filter deleted first
	dbc := db.conn.Table(models.Account{}.TableName())

	if err := dbc.Where("id = ?", id).Update("IsDisable", isDisable).Error; err != nil {
		return err
	}

	return nil
}
