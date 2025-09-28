package storage

import (
	"context"
	"github.com/coderconquerer/social-todo/module/todotaskreaction/entity"
)

func (db *MySQLConnection) CreateReaction(c context.Context, reaction entity.Reaction) error {

	dbc := db.conn.Table(entity.Reaction{}.TableName())

	if err := dbc.Create(reaction).Error; err != nil {
		return err
	}

	return nil
}
