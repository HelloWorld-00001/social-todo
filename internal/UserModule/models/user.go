package models

func (User) TableName() string {
	return "account"
}

type User struct {
	UserID   int    `gorm:"column:User_Id;primaryKey;autoIncrement" json:"user_id"`
	Email    string `gorm:"column:Email;unique" json:"email"`
	Phone    string `gorm:"column:Phone;unique" json:"phone"`
	Username string `gorm:"column:Username;unique" json:"username"`
	Name     string `gorm:"column:Name" json:"name"`
	Avatar   *int   `gorm:"column:avatar" json:"avatar"`
}

func (Account) TableName() string {
	return "account"
}

type Account struct {
	AccountID int    `gorm:"column:Account_Id;primaryKey;autoIncrement" json:"account_id"`
	Password  string `gorm:"column:Password" json:"password"`
	Salt      string `gorm:"column:Salt" json:"salt"`
	Username  string `gorm:"column:Username;unique" json:"username"`
}
