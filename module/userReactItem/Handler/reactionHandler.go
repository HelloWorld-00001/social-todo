package Handler

import "github.com/coderconquerer/social-todo/module/userReactItem/BusinessUseCases"

type ReactionHandler struct {
	UnreactTodoItemBz *BusinessUseCases.UnreactTodoItemLogic
	ReactTodoItemBz   *BusinessUseCases.ReactTodoItemLogic
	ListReactedUserBz *BusinessUseCases.GetListReactedUsersLogic
}

func NewReactionHandler(reactTodoItemBz *BusinessUseCases.ReactTodoItemLogic,
	unreactTodoItemBz *BusinessUseCases.UnreactTodoItemLogic,
	listReactedUserBz *BusinessUseCases.GetListReactedUsersLogic) *ReactionHandler {
	return &ReactionHandler{
		ReactTodoItemBz:   reactTodoItemBz,
		UnreactTodoItemBz: unreactTodoItemBz,
		ListReactedUserBz: listReactedUserBz,
	}
}
