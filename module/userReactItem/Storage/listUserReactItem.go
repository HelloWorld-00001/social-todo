package Storage

import (
	"context"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/common/helper"
	models2 "github.com/coderconquerer/social-todo/module/user/models"
	"github.com/coderconquerer/social-todo/module/userReactItem/models"
	"github.com/gin-gonic/gin"
	"go.opencensus.io/trace"
)

func (db *MySQLConnection) GetReactedUsers(c *gin.Context, todoId int, pagination *common.Pagination) ([]models2.SimpleUser, error) {
	var reactions []models.Reaction
	// filter deleted first
	dbc := db.conn.Table(models.Reaction{}.TableName()).Where("TodoId = ?", todoId)

	if err := dbc.Select("UserId").Count(&pagination.Total).Error; err != nil {
		return nil, err
	}

	if pagination.Cursor != "" {
		createdAt, err := helper.DecodeBase64URLToTime(pagination.Cursor)
		if err != nil {
			return nil, err
		}

		dbc.Where("CreatedAt < ?", createdAt.Format(helper.MySqlTimeLayout))
	} else {
		dbc.Offset((pagination.Page - 1) * pagination.Limit)
	}

	if err := dbc.Select("*").
		Order("CreatedAt desc").
		Limit(pagination.Limit).
		Preload("ReactedUser").
		Find(&reactions).Error; err != nil {
		return nil, err
	}

	size := len(reactions)
	users := make([]models2.SimpleUser, size)

	for i := range users {
		users[i] = *reactions[i].ReactedUser
		users[i].CreatedAt = nil
		users[i].UpdatedAt = nil
		users[i].React = reactions[i].React.String()
		users[i].ReactedAt = reactions[i].CreatedAt
	}

	if size == pagination.Limit {
		pagination.NextCursor = helper.EncodeTimeToBase64URL(reactions[size-1].CreatedAt)
	}

	return users, nil
}

func (db *MySQLConnection) GetReactedTodo(c context.Context, todoIds []int) (map[int]int, error) {
	_, span := trace.StartSpan(c, "todo_react.storage.GetReactedTodo")
	defer span.End()
	dbc := db.conn.Table(models.Reaction{}.TableName())
	type AggReact struct {
		LikeCount int
		TodoId    int
	}

	var reactions []AggReact
	if err := dbc.Select("TodoId, COUNT(TodoId) AS LikeCount").
		Where("TodoId IN ?", todoIds).
		Group("TodoId").
		Scan(&reactions).Error; err != nil {
		return nil, err
	}
	mp := make(map[int]int, len(reactions))
	for _, val := range reactions {
		mp[val.TodoId] = val.LikeCount
	}

	return mp, nil
}
