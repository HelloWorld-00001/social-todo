package Storage

import (
	"github.com/coderconquerer/social-todo/common"
	models2 "github.com/coderconquerer/social-todo/module/user/models"
	"github.com/coderconquerer/social-todo/module/userReactItem/models"
	"github.com/gin-gonic/gin"
)

func (db *MySQLConnection) GetReactedUsers(c *gin.Context, todoId int, pagination *common.Pagination) ([]models2.SimpleUser, error) {
	var reactions []models.Reaction
	// filter deleted first
	dbc := db.conn.Table(models.Reaction{}.TableName()).Where("TodoId = ?", todoId)

	if err := dbc.Select("UserId").Count(&pagination.Total).Error; err != nil {
		return nil, err
	}

	if err := dbc.Select("*").
		Order("CreatedAt desc").
		Offset((pagination.Page - 1) * pagination.Limit).
		Limit(pagination.Limit).
		Preload("ReactedUser").
		Find(&reactions).Error; err != nil {
		return nil, err
	}

	users := make([]models2.SimpleUser, len(reactions))
	for i := range users {
		users[i] = *reactions[i].ReactedUser
		users[i].React = reactions[i].React.String()
	}
	return users, nil
}
