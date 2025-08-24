package Storage

import (
	"context"
	"errors"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/account/models"
	models2 "github.com/coderconquerer/social-todo/module/user/models"
	"github.com/gin-gonic/gin"
	"go.opencensus.io/trace"
	"gorm.io/gorm"
)

func (db *MySQLConnection) FindAccountByUsername(c *gin.Context, username string) (*models.Account, error) {
	_, span := trace.StartSpan(c, "user.storage.FindAccountByUsername")
	defer span.End()
	// filter deleted first
	var account models.Account
	dbc := db.conn.Table(models.Account{}.TableName())
	if err := dbc.Preload(models2.User{}.TableName()).Where("Username = ?", username).Take(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, common.NewDatabaseError(err)
	}

	return &account, nil
}

func (db *MySQLConnection) FindAccount(c context.Context, conditions map[string]interface{}) (*models.Account, error) {
	_, span := trace.StartSpan(c, "user.storage.FindAccount")
	defer span.End()

	var account models.Account
	dbc := db.conn.Table(models.Account{}.TableName())

	if err := dbc.Preload(models2.User{}.TableName()).Where(conditions).Take(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}
