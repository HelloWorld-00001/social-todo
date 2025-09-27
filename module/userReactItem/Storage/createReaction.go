package Storage

import (
	"context"
	"github.com/coderconquerer/social-todo/module/userReactItem/models"
)

func (db *MySQLConnection) CreateReaction(c context.Context, reaction models.Reaction) error {
	// filter deleted first
	dbc := db.conn.Table(models.Reaction{}.TableName())

	if err := dbc.Create(reaction).Error; err != nil {
		return err
	}

	return nil
}
