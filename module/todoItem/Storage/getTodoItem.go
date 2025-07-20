package Storage

import (
	"github.com/coderconquerer/social-todo/module/todoItem/models"
	"github.com/gin-gonic/gin"
)

func (db *MySQLConnection) GetTodoItemDetailById(c *gin.Context, id int) (*models.Todo, error) {
	var todo models.Todo
	// filter deleted first
	dbc := db.conn.Table(models.Todo{}.TableName())

	// Todo: if account's role = admin -> return item despite of being deleted or not,
	// otherwise, return not found if record is deleted
	dbc = dbc.Where("Deleted_Date is null")

	if err := dbc.Where("Id = ?", id).Take(&todo).Error; err != nil {
		return nil, err
	}

	return &todo, nil
}
