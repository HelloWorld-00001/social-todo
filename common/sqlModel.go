package common

import "time"

type SqlModel struct {
	Id        int        `json:"-" gorm:"column:Id;primaryKey;autoIncrement;"`
	MarkupId  *Uid       `json:"id" gorm:"-"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:CreatedAt;"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:UpdatedAt;"`
}

func (s *SqlModel) MakeMarkupId(dbType Entity, shardId int) {
	uid := NewUid(uint32(s.Id), int(dbType), 1)
	s.MarkupId = uid
}
