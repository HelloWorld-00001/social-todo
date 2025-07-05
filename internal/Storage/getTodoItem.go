package Storage

import (
	"errors"
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/coderconquerer/go-login-app/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (db *MySQLConnection) GetTodoItemDetailById(c *gin.Context, id int) (*models.Todo, error) {
	var todo models.Todo
	// filter deleted first
	dbc := db.conn.Table(models.Todo{}.TableName())

	// Todo: if user's role = admin -> return item despite of being deleted or not,
	// otherwise, return not found if record is deleted
	dbc = dbc.Where("Deleted_Date is null")

	if err := dbc.Where("Id = ?", id).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, common.NewDatabaseError(err)
	}

	return &todo, nil
}
