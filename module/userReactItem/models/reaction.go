package models

import (
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/user/models"
	"time"
)

func (Reaction) TableName() string {
	return "Reaction"
}

type Reaction struct {
	UserId    int          `gorm:"column:UserId"`
	TodoId    int          `gorm:"column:TodoId"`
	React     common.React `gorm:"column:React;type:enum('Like','Dislike','Love','Angry','Wow')"`
	CreatedAt time.Time    `gorm:"column:CreatedAt"`

	ReactedUser *models.SimpleUser `gorm:"foreignKey:UserId" json:"reacted_users"`
}

type ReactionInput struct {
	TodoId string `form:"todo_id" binding:"required"`
	React  string `form:"reaction" binding:"required"`
}
