package common

type SqlModel struct {
	Id       int  `json:"-" gorm:"column:Id;primaryKey;autoIncrement;"`
	MarkupId *Uid `json:"id" gorm:"-"`
}

func (s *SqlModel) MakeMarkupId(dbType Entity, shardId int) {
	uid := NewUid(uint32(s.Id), int(dbType), 1)
	s.MarkupId = uid
}
