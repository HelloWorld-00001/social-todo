package models

import (
	common "github.com/coderconquerer/social-todo/common"
	"time"
)

func (User) TableName() string {
	return "User"
}

type User struct {
	common.SqlModel
	Email    string        `gorm:"column:Email;unique" json:"email"`
	Phone    string        `gorm:"column:Phone;unique" json:"phone"`
	Username string        `gorm:"column:Username;unique" json:"username"`
	Name     string        `gorm:"column:Name" json:"name"`
	Avatar   *common.Image `gorm:"column:avatar" json:"avatar"`
}

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

func (u *User) CreateMarkupId() {
	u.MakeMarkupId(common.UserEntity, 1)
}

func (u *SimpleUser) CreateMarkupId() {
	u.MakeMarkupId(common.UserEntity, 1)
}
