package Storage

import (
	"errors"
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/coderconquerer/go-login-app/internal/user/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (db *MySQLConnection) FindUserById(c *gin.Context, id int) (*models.User, error) {
	// filter deleted first
	var user models.User
	dbc := db.conn.Table(models.User{}.TableName())
	if err := dbc.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, common.NewDatabaseError(err)
	}

	return &user, nil
}

func (db *MySQLConnection) FindUser(c *gin.Context, conditions map[string]interface{}) (*models.User, error) {
	var user models.User
	dbc := db.conn.Table(models.User{}.TableName())

	if err := dbc.Where(conditions).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
