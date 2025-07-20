package Storage

import (
	"github.com/coderconquerer/social-todo/common"
	model "github.com/coderconquerer/social-todo/module/account/models"
	"github.com/gin-gonic/gin"
)

func (db *MySQLConnection) LoginAccount(c *gin.Context, account model.AccountLogin) (int, error) {
	// filter deleted first
	dbc := db.conn.Table(model.Account{}.TableName())
	id := common.InvalidID
	if err := dbc.Select("Id").Where("Username = ? & Password = ?", account.Username, account.Password).First(id).Error; err != nil {
		return common.InvalidID, err
	}

	return id, nil
}
