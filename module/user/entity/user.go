package entity

import (
	"github.com/coderconquerer/social-todo/common"
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

func (u *User) CreateMarkupId() {
	u.MakeMarkupId(common.UserEntity, 1)
}
