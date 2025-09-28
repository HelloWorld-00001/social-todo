package entity

type AccountLogin struct {
	Password string `gorm:"column:Password" json:"password"`
	Username string `gorm:"column:Username;unique" json:"username"`
}

type AccountRegister struct {
	Password string `json:"password"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
