package models

func (Account) TableName() string {
	return "Account"
}

type Account struct {
	AccountID int    `gorm:"column:Id;primaryKey;autoIncrement" json:"-"`
	Password  string `gorm:"column:Password" json:"password"`
	Salt      string `gorm:"column:Salt" json:"-"`
	Username  string `gorm:"column:Username;unique" json:"username"`
	Role      string `gorm:"column:Role" json:"role"`
	IsDisable bool   `gorm:"column:IsDisable" json:"IsDisable"`
}

type AccountLogin struct {
	Password string `gorm:"column:Password" json:"password"`
	Username string `gorm:"column:Username;unique" json:"username"`
}
