package models

import (
	common2 "github.com/coderconquerer/go-login-app/common"
)

func (User) TableName() string {
	return "User"
}

type User struct {
	common2.SqlModel
	Email    string         `gorm:"column:Email;unique" json:"email"`
	Phone    string         `gorm:"column:Phone;unique" json:"phone"`
	Username string         `gorm:"column:Username;unique" json:"username"`
	Name     string         `gorm:"column:Name" json:"name"`
	Avatar   *common2.Image `gorm:"column:avatar" json:"avatar"`
}

func (u *User) MarkupId() {
	u.MakeMarkupId(common2.UserEntity, 1)
}
