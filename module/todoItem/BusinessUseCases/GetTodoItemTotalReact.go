package BusinessUseCases

import (
	"github.com/gin-gonic/gin"
)

type ReactionCountStorage interface {
	GetReactedTodo(c *gin.Context, todoIds []int) (map[int]int, error)
}

type GetTodoItemTotalReactLogic struct {
	store ReactionCountStorage
}

func GetNewGetTodoItemTotalReactLogic(store ReactionCountStorage) *GetTodoItemTotalReactLogic {
	return &GetTodoItemTotalReactLogic{store: store}
}

func (bz *GetTodoItemTotalReactLogic) GetTodoItemTotalReact(c *gin.Context, todoIds []int) (map[int]int, error) {
	data, err := bz.store.GetReactedTodo(c, todoIds)
	if err != nil {
		return nil, err
	}
	return data, nil
}
