package common

type Responser struct {
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination,omitempty"`
	Filter     interface{} `json:"filter,omitempty"`
}

func SimpleResponse(data interface{}) *Responser {
	return &Responser{Data: data}
}

func StandardResponse(data interface{}, pagination Pagination, filter Filter) *Responser {
	return &Responser{Data: data, Pagination: pagination, Filter: filter}
}

func StandardResponseWithoutFilter(data interface{}, pagination Pagination) *Responser {
	return &Responser{Data: data, Pagination: pagination}
}
