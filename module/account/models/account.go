package models

import models2 "github.com/coderconquerer/social-todo/module/user/models"

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

	User *models2.User `gorm:"foreignKey:Username;references:Username" json:"user"`
}

type AccountLogin struct {
	Password string `gorm:"column:Password" json:"password"`
	Username string `gorm:"column:Username;unique" json:"username"`
}
