package entity

type AccountLogin struct {
	Password string `gorm:"column:Password" json:"password"`
	Username string `gorm:"column:Username;unique" json:"username"`
}
