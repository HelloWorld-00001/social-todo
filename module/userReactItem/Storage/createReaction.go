package Storage

import (
	"github.com/coderconquerer/social-todo/module/userReactItem/models"
	"github.com/gin-gonic/gin"
)

func (db *MySQLConnection) CreateReaction(c *gin.Context, reaction models.Reaction) error {
	// filter deleted first
	dbc := db.conn.Table(models.Reaction{}.TableName())

	if err := dbc.Create(reaction).Error; err != nil {
		return err
	}

	return nil
}
