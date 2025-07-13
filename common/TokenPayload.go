package common

type TokenPayload struct {
	UserId int    `json:"userId"`
	Role   string `json:"role"`
}

func (tp TokenPayload) GetUserId() int {
	return tp.UserId
}

func (tp TokenPayload) GetRole() string {
	return tp.Role
}
