package common

type TokenPayload struct {
	AccountId int    `json:"account_id"`
	UserId    int    `json:"userId"`
	Role      string `json:"role"`
}

func (tp TokenPayload) GetUserId() int {
	return tp.UserId
}

func (tp TokenPayload) GetAccountId() int {
	return tp.AccountId
}

func (tp TokenPayload) GetRole() string {
	return tp.Role
}
