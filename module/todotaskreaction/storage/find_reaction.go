package storage

import (
	"context"
	"github.com/coderconquerer/social-todo/module/todotaskreaction/entity"
)

func (db *MySQLConnection) FindReaction(c context.Context, userId, todoId int) (*entity.Reaction, error) {
	react := &entity.Reaction{}
	dbc := db.conn.Table(entity.Reaction{}.TableName())

	if err := dbc.Where("TodoId = ? and UserId = ?", todoId, userId).Take(react).Error; err != nil {
		return nil, err
	}

	return react, nil
}
