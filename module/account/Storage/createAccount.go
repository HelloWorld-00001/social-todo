package Storage

import (
	"github.com/coderconquerer/go-login-app/module/account/models"
	"github.com/gin-gonic/gin"
)

func (db *MySQLConnection) CreateAccount(c *gin.Context, acc *models.Account) error {
	// filter deleted first
	dbc := db.conn.Table(models.Account{}.TableName())

	if err := dbc.Create(acc).Error; err != nil {
		return err
	}

	return nil
}
