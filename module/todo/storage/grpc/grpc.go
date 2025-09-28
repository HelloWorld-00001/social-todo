package restapi

import (
	"context"
	"github.com/coderconquerer/social-todo/common/helper"
	"github.com/coderconquerer/social-todo/grpc/contract"
	"google.golang.org/grpc"
)

type rpcClient struct {
	client contract.ItemReactServiceClient
}

func NewRpcClient(client contract.ItemReactServiceClient) *rpcClient {
	return &rpcClient{
		client: client,
	}
}

func (rc *rpcClient) GetTodoTotalReact(c context.Context, todoIds []int) (map[int]int, error) {

	ids := helper.ListIntToInt32(todoIds)

	res, err := rc.client.GetTotalReactByIds(c, &contract.GetTotalReactByIdsRequest{Ids: ids})

	if err != nil {
		return nil, err
	}

	result := helper.MapInt32ToInt(res.Result)
	return result, nil
}

func (rc *rpcClient) GetTotalReactByIds(c context.Context, todoIds *contract.GetTotalReactByIdsRequest, opts ...grpc.CallOption) (*contract.GetTotalReactByIdsResponse, error) {
	//ids := helper.ListInt32ToInt(todoIds.Ids)
	//
	//res, err := rc.client.GetTotalReactByIds(c, &contract.GetTotalReactByIdsRequest{Ids: ids})
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//result := helper.MapInt32ToInt(res.Result)
	// todo-implement later
	return nil, nil
}
