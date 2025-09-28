package entity

import (
	"github.com/coderconquerer/social-todo/common"
	"time"
)

func (SimpleUser) TableName() string {
	return "User"
}

type SimpleUser struct {
	common.SqlModel
	Email     string    `gorm:"column:Email;unique" json:"email"`
	Name      string    `gorm:"column:Name" json:"name"`
	React     string    `gorm:"-" json:"react"`
	ReactedAt time.Time `gorm:"-" json:"reacted_at"`
}

func (u *SimpleUser) CreateMarkupId() {
	u.MakeMarkupId(common.UserEntity, 1)
}
