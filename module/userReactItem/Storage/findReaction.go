package Storage

import (
	"github.com/coderconquerer/social-todo/module/userReactItem/models"
	"github.com/gin-gonic/gin"
)

func (db *MySQLConnection) FindReaction(c *gin.Context, userId, todoId int) (*models.Reaction, error) {
	react := &models.Reaction{}
	dbc := db.conn.Table(models.Reaction{}.TableName())

	if err := dbc.Where("TodoId = ? and UserId = ?", todoId, userId).Take(react).Error; err != nil {
		return nil, err
	}

	return react, nil
}
