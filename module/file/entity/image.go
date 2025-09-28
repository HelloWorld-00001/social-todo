package entity

func (User) TableName() string {
	return "User"
}

type User struct {
	UserID   int    `gorm:"column:Id;primaryKey;autoIncrement" json:"user_id"`
	Email    string `gorm:"column:Email;unique" json:"email"`
	Phone    string `gorm:"column:Phone;unique" json:"phone"`
	Username string `gorm:"column:Username;unique" json:"username"`
	Name     string `gorm:"column:Name" json:"name"`
	Avatar   *int   `gorm:"column:avatar" json:"avatar"`
}
