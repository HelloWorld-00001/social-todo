package Storage

import (
	"errors"
	"github.com/coderconquerer/go-login-app/internal/TodoItem/models"
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (db *MySQLConnection) GetUserProfile(c *gin.Context, id int) (*models.User, error) {
	var user models.User
	// filter deleted first
	dbc := db.conn.Table(models.User{}.TableName())

	// Todo: if account's role = admin -> return item despite of being deleted or not,
	// otherwise, return not found if record is deleted
	dbc = dbc.Where("Deleted_Date is null")

	if err := dbc.Where("Id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, common.NewDatabaseError(err)
	}

	return &user, nil
}
