package Storage

import (
	"github.com/coderconquerer/social-todo/module/userReactItem/models"
	"github.com/gin-gonic/gin"
)

func (db *MySQLConnection) DeleteReaction(c *gin.Context, userId, todoId int) error {
	dbc := db.conn.Table(models.Reaction{}.TableName())

	if err := dbc.Where("TodoId = ? AND UserId = ?", todoId, userId).
		Delete(&models.Reaction{}).Error; err != nil {
		return err
	}

	return nil
}
