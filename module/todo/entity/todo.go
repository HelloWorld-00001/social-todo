package entity

import (
	"github.com/coderconquerer/social-todo/common"
	"time"
)

func (Todo) TableName() string {
	return TodoTableName
}

var (
	TodoTableName = "Todo"
)

type Todo struct {
	common.SqlModel
	Description string     `gorm:"column:Description;type:text;" json:"description"`
	Status      string     `gorm:"column:Status" json:"status"`
	UpdateTime  time.Time  `gorm:"column:UpdateTime" json:"update_time"`
	CreateTime  time.Time  `gorm:"column:CreateTime" json:"create_time"`
	Deadline    time.Time  `gorm:"column:Deadline" json:"deadline"`
	DeletedDate *time.Time `gorm:"column:Deleted_Date" json:"deleted_date,omitempty"`
	Label       string     `gorm:"column:Label" json:"label"`
	TagColor    string     `gorm:"column:TagColor" json:"tag_color"`
	Workspace   string     `gorm:"column:workspace" json:"workspace"`
	TotalReact  int        `gorm:"column:TotalReact" json:"total_react"`
	CreateBy    int        `gorm:"column:Create_By" json:"create_by"`
	Assignee    int        `gorm:"column:Assignee" json:"assignee"`
	Title       string     `gorm:"column:Title" json:"title"`

	//Creator  *User `gorm:"foreignKey:CreateBy;references:UserId" json:"creator,omitempty"`
	//Assigned *User `gorm:"foreignKey:Assignee;references:UserId" json:"assigned,omitempty"`
}

func (t *Todo) CreateMarkupId() {
	t.MakeMarkupId(common.TodoEntity, 1)
}
