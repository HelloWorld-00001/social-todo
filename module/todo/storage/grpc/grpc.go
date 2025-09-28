package restapi

import (
	"context"
	"fmt"
	"github.com/coderconquerer/social-todo/common/helper"
	"github.com/coderconquerer/social-todo/grpc/contract"
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

	fmt.Println(ids)
	fmt.Println("hehe it ok")
	res, err := rc.client.GetTotalReactByIds(c, &contract.GetTotalReactByIdsRequest{Ids: ids})

	if err != nil {
		return nil, err
	}

	result := helper.MapInt32ToInt(res.Result)
	return result, nil
}
