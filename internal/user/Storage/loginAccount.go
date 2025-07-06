package Storage

import (
	"github.com/coderconquerer/go-login-app/internal/TodoItem/models"
	model "github.com/coderconquerer/go-login-app/internal/account/models"
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/gin-gonic/gin"
)

func (db *MySQLConnection) LoginAccount(c *gin.Context, account model.AccountLogin) (int, error) {
	// filter deleted first
	dbc := db.conn.Table(models.Account{}.TableName())
	id := common.InvalidID
	if err := dbc.Select("Id").Where("Username = ? & Password = ?", account.Username, account.Password).First(id).Error; err != nil {
		return common.InvalidID, err
	}

	return id, nil
}
