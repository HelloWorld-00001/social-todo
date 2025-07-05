package Storage

import (
	"github.com/coderconquerer/go-login-app/internal/models"
	"github.com/gin-gonic/gin"
	"time"
)

func (db *MySQLConnection) CreateTodoItem(c *gin.Context, todo *models.TodoCreation) error {
	// filter deleted first
	dbc := db.conn.Table(models.Todo{}.TableName())
	currentTime := time.Now()
	todo.CreateTime = currentTime
	todo.UpdateTime = currentTime

	if err := dbc.Create(todo).Error; err != nil {
		return err
	}

	return nil
}
