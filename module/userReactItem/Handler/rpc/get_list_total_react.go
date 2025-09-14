package rpc

import (
	"context"
	"github.com/coderconquerer/social-todo/common/helper"
	"github.com/coderconquerer/social-todo/grpc/contract"
	"github.com/coderconquerer/social-todo/module/userReactItem/BusinessUseCases"
)

type rpcService struct {
	listReactedUserBz *BusinessUseCases.GetListReactedUsersLogic
}

func NewRpcService(
	listReactedUserBz *BusinessUseCases.GetListReactedUsersLogic) contract.ItemReactServiceServer {
	return &rpcService{
		listReactedUserBz: listReactedUserBz,
	}
}

func (rh rpcService) GetTotalReactByIds(ctx context.Context, request *contract.GetTotalReactByIdsRequest) (*contract.GetTotalReactByIdsResponse, error) {
	ids := helper.ListInt32ToInt(request.Ids)

	res, err := rh.listReactedUserBz.GetTodoItemTotalReact(ctx, ids)
	if err != nil {
		return nil, err
	}

	result := helper.MapIntToInt32(res)
	return &contract.GetTotalReactByIdsResponse{Result: result}, nil
}
