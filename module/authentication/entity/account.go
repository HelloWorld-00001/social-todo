package entity

import userEntity "github.com/coderconquerer/social-todo/module/user/entity"

func (Account) TableName() string {
	return "Account"
}

type Account struct {
	Id        int    `gorm:"column:Id;primaryKey;autoIncrement" json:"-"`
	Password  string `gorm:"column:Password" json:"password"`
	Salt      string `gorm:"column:Salt" json:"-"`
	Username  string `gorm:"column:Username;unique" json:"username"`
	Role      string `gorm:"column:Role" json:"role"`
	IsDisable bool   `gorm:"column:IsDisable" json:"IsDisable"`

	User *userEntity.User `gorm:"foreignKey:Username;references:Username" json:"user"`
}
