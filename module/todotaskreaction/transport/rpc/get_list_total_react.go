package rpc

import (
	"context"
	"github.com/coderconquerer/social-todo/common/helper"
	"github.com/coderconquerer/social-todo/grpc/contract"
)

type ReactionBusinessRPC interface {
	GetTodoItemTotalReact(c context.Context, todoIds []int) (map[int]int, error)
}
type rpcService struct {
	business ReactionBusinessRPC
}

func NewRpcService(business ReactionBusinessRPC) contract.ItemReactServiceServer {
	return &rpcService{
		business,
	}
}

func (rh rpcService) GetTotalReactByIds(ctx context.Context, request *contract.GetTotalReactByIdsRequest) (*contract.GetTotalReactByIdsResponse, error) {
	ids := helper.ListInt32ToInt(request.Ids)

	res, err := rh.business.GetTodoItemTotalReact(ctx, ids)
	if err != nil {
		return nil, err
	}
	result := helper.MapIntToInt32(res)
	return &contract.GetTotalReactByIdsResponse{Result: result}, nil
}
