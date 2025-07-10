package models

func (Account) TableName() string {
	return "account"
}

type Account struct {
	AccountID int    `gorm:"column:Account_Id;primaryKey;autoIncrement" json:"account_id"`
	Password  string `gorm:"column:Password" json:"password"`
	Salt      string `gorm:"column:Salt" json:"salt"`
	Username  string `gorm:"column:Username;unique" json:"username"`
}

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

	AvatarMetadata *Metadata  `gorm:"foreignKey:Avatar;references:MetadataID" json:"avatar_metadata,omitempty"`
	CreatedTodos   []Todo     `gorm:"foreignKey:CreateBy" json:"created_todos,omitempty"`
	AssignedTodos  []Todo     `gorm:"foreignKey:Assignee" json:"assigned_todos,omitempty"`
	Comments       []Comment  `gorm:"foreignKey:Id" json:"comments,omitempty"`
	Reactions      []Reaction `gorm:"foreignKey:Id" json:"reactions,omitempty"`
}

func (Metadata) TableName() string {
	return "metadata"
}

type Metadata struct {
	MetadataID    int    `gorm:"column:metadata_id;primaryKey;autoIncrement" json:"metadata_id"`
	FileName      string `gorm:"column:FileName" json:"file_name"`
	FileExtension string `gorm:"column:FileExtension" json:"file_extension"`
	FileSize      int    `gorm:"column:FileSize" json:"file_size"`
	Height        int    `gorm:"column:Height" json:"height"`
	Width         int    `gorm:"column:Width" json:"width"`
	URL           string `gorm:"column:URL" json:"url"`
	TodoID        *int   `gorm:"column:Todo_id" json:"todo_id,omitempty"`

	Todo *Todo `gorm:"foreignKey:Id" json:"todo,omitempty"`
}

func (Comment) TableName() string {
	return "comment"
}

type Comment struct {
	TodoID  int    `gorm:"column:Todo_Id;primaryKey" json:"todo_id"`
	UserID  int    `gorm:"column:User_Id;primaryKey" json:"user_id"`
	Content string `gorm:"column:Content;type:text;" json:"content"`

	Todo *Todo `gorm:"foreignKey:Id" json:"todo,omitempty"`
	User *User `gorm:"foreignKey:Id" json:"account,omitempty"`
}

func (Reaction) TableName() string {
	return "reaction"
}

type Reaction struct {
	TodoID int    `gorm:"column:Todo_Id;primaryKey" json:"todo_id"`
	UserID int    `gorm:"column:User_Id;primaryKey" json:"user_id"`
	React  string `gorm:"column:React;type:enum('Like','Dislike','Love','Angry','Wow')" json:"react"`

	Todo *Todo `gorm:"foreignKey:Id" json:"todo,omitempty"`
	User *User `gorm:"foreignKey:Id" json:"account,omitempty"`
}
