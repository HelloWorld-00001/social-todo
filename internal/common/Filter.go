package common

import (
	"github.com/coderconquerer/go-login-app/internal/TodoItem/models"
	"time"
)

type Filter struct {
	Description string     `gorm:"column:Description;type:text;" json:"description"`
	Status      string     `gorm:"column:Status" json:"status"`
	UpdateTime  time.Time  `gorm:"column:UpdateTime" json:"update_time"`
	CreateTime  time.Time  `gorm:"column:CreateTime" json:"create_time"`
	Deadline    time.Time  `gorm:"column:Deadline" json:"deadline"`
	DeletedDate *time.Time `gorm:"column:DeletedDate" json:"deleted_date,omitempty"`
	Label       string     `gorm:"column:Label" json:"label"`
	TagColor    string     `gorm:"column:TagColor" json:"tag_color"`
	Workspace   string     `gorm:"column:workspace" json:"workspace"`
	CreateBy    int        `gorm:"column:Create_By" json:"create_by"`
	Assignee    int        `gorm:"column:Assignee" json:"assignee"`
	Title       string     `gorm:"column:Title" json:"title"`

	Creator  *models.User `gorm:"foreignKey:CreateBy;references:UserID" json:"creator,omitempty"`
	Assigned *models.User `gorm:"foreignKey:Assignee;references:UserID" json:"assigned,omitempty"`
}
