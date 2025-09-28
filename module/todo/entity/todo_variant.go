package entity

import (
	common "github.com/coderconquerer/social-todo/common"
	"time"
)

type UpdateTodo struct {
	common.SqlModel
	Description string    `gorm:"column:Description;type:text;" json:"description"`
	Status      string    `gorm:"column:Status" json:"status"`
	UpdateTime  time.Time `gorm:"column:UpdateTime" json:"-"`
	Deadline    time.Time `gorm:"column:Deadline" json:"deadline"`
	Label       string    `gorm:"column:Label" json:"label"`
	TagColor    string    `gorm:"column:TagColor" json:"tag_color"`
	Workspace   string    `gorm:"column:workspace" json:"workspace"`
	Assignee    int       `gorm:"column:Assignee" json:"assignee"`
	Title       string    `gorm:"column:Title" json:"title"`
}

type TodoCreation struct {
	TodoID      int       `gorm:"column:Id;primaryKey;autoIncrement" json:"-"`
	Description string    `gorm:"column:Description;type:text;" json:"description"`
	Status      string    `gorm:"column:Status" json:"status"`
	UpdateTime  time.Time `gorm:"column:UpdateTime" json:"-"`
	CreateTime  time.Time `gorm:"column:CreateTime" json:"-"`
	Deadline    time.Time `gorm:"column:Deadline" json:"deadline"`
	Label       string    `gorm:"column:Label" json:"label"`
	TagColor    string    `gorm:"column:TagColor" json:"tag_color"`
	Workspace   string    `gorm:"column:workspace" json:"workspace"`
	Assignee    int       `gorm:"column:Assignee" json:"assignee"`
	Title       string    `gorm:"column:Title" json:"title"`
}
