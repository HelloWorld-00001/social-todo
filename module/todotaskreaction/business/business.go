package business

import (
	"context"
	"errors"
	"github.com/coderconquerer/social-todo/cmd/registerservice"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/todotaskreaction/entity"
	userEntity "github.com/coderconquerer/social-todo/module/user/entity"
	"github.com/coderconquerer/social-todo/pubsub"
	"gorm.io/gorm"
	"log"
)

type ReactionStorage interface {
	CreateReaction(c context.Context, reaction entity.Reaction) error
	GetReactedUsers(c context.Context, todoId int, pagination *common.Pagination) ([]userEntity.SimpleUser, error)
	GetReactedTodo(c context.Context, todoIds []int) (map[int]int, error)
	FindReaction(c context.Context, userId, todoId int) (*entity.Reaction, error)
	DeleteReaction(c context.Context, userId, todoId int) error
}

type ReactionBusiness interface {
	ReactTodoItem(c context.Context, reaction entity.Reaction) error
	GetListReactedUsers(c context.Context, todoId int, pagination *common.Pagination) ([]userEntity.SimpleUser, error)
	GetTodoItemTotalReact(c context.Context, todoIds []int) (map[int]int, error)
	UnreactTodoItem(c context.Context, userId, todoId int) error
}
type reactionBusiness struct {
	store    ReactionStorage
	ps       pubsub.PubSub
	rabbitPs registerservice.RbPublisher
}

func NewReactionBusiness(store ReactionStorage, ps pubsub.PubSub, rabbitPs registerservice.RbPublisher) ReactionBusiness {
	return &reactionBusiness{store, ps, rabbitPs}
}

func (bz *reactionBusiness) ReactTodoItem(c context.Context, reaction entity.Reaction) error {
	err := bz.store.CreateReaction(c, reaction)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NotFound.WithError(entity.ErrCannotFindReaction).WithRootCause(err)
		}
		return common.InternalServerError.WithError(common.ErrUnhandleError).WithRootCause(err)
	}

	errPs := bz.ps.Publish(c, common.TopicIncreaseTotalReact, pubsub.NewMessage(reaction))
	if errPs != nil {
		log.Println(errPs)
	}
	return nil
}

func (bz *reactionBusiness) GetListReactedUsers(c context.Context, todoId int, pagination *common.Pagination) ([]userEntity.SimpleUser, error) {
	data, err := bz.store.GetReactedUsers(c, todoId, pagination)
	if err != nil {
		return nil, common.NotFound.WithError(entity.ErrCannotFindReaction).WithRootCause(err)
	}
	return data, nil
}

func (bz *reactionBusiness) GetTodoItemTotalReact(c context.Context, todoIds []int) (map[int]int, error) {
	data, err := bz.store.GetReactedTodo(c, todoIds)
	if err != nil {
		return nil, common.InternalServerError.WithError(errors.New("error when trying to get reacted todo task")).WithRootCause(err)
	}
	return data, nil
}

func (bz *reactionBusiness) UnreactTodoItem(c context.Context, userId, todoId int) error {
	react, err := bz.store.FindReaction(c, userId, todoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NotFound.WithError(entity.ErrCannotFindReaction).WithRootCause(err)
		}
		return common.InternalServerError.WithError(common.ErrUnhandleError).WithRootCause(err)
	}

	errDel := bz.store.DeleteReaction(c, react.UserId, react.TodoId)
	if errDel != nil {
		return common.InternalServerError.WithError(entity.ErrCannotDeleteReaction).WithRootCause(errDel)
	}

	var rct entity.Reaction
	rct.UserId = userId
	rct.TodoId = todoId

	errPs := bz.ps.Publish(c, common.TopicDecreaseTotalReact, pubsub.NewMessage(rct))

	if errPs != nil {
		log.Println(errPs)
	}
	return nil
}
