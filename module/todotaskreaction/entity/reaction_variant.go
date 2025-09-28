package entity

type ReactionInput struct {
	TodoId string `form:"todo_id" binding:"required"`
	React  string `form:"reaction" binding:"required"`
}
