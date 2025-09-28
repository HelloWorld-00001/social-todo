package storage

import (
	"context"
	"github.com/coderconquerer/social-todo/module/todotaskreaction/entity"
)

func (db *MySQLConnection) DeleteReaction(c context.Context, userId, todoId int) error {
	dbc := db.conn.Table(entity.Reaction{}.TableName())

	if err := dbc.Where("TodoId = ? AND UserId = ?", todoId, userId).
		Delete(&entity.Reaction{}).Error; err != nil {
		return err
	}

	return nil
}
