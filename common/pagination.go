package common

type Pagination struct {
	Total      int64  `json:"total" form:"-"`
	Page       int    `json:"page,omitempty" form:"page"`
	Limit      int    `json:"limit,omitempty" form:"limit"`
	Cursor     string `json:"cursor" form:"cursor"`
	NextCursor string `json:"next_cursor" form:"next_cursor"`
}

func (p *Pagination) Process() {
	if p.Limit <= 0 {
		p.Limit = 5
	}

	if p.Limit > 50 {
		p.Limit = 50
	}

	if p.Page < 1 {
		p.Page = 1
	}

	if p.Page > p.Limit {
		p.Page = p.Limit
	}
}
