package entity

import (
	"github.com/coderconquerer/social-todo/common"
	userEntity "github.com/coderconquerer/social-todo/module/user/entity"
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

	ReactedUser *userEntity.SimpleUser `gorm:"foreignKey:UserId" json:"reacted_users"`
}

func (r Reaction) GetTodoId() int {
	return r.TodoId
}

func (r Reaction) GetUserId() int {
	return r.UserId
}

func (r Reaction) GetReaction() string {
	return r.React.String()
}
