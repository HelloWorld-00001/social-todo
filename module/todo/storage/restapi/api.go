package restapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
)

type todoReactService struct {
	serviceUrl string
	client     *resty.Client
}

func NewTodoReactService(serviceUrl string, client *resty.Client) *todoReactService {
	return &todoReactService{
		serviceUrl: serviceUrl,
		client:     client,
	}
}

func (ts *todoReactService) GetTodoTotalReact(c context.Context, todoIds []int) (map[int]int, error) {
	type requestBody struct {
		Ids []int `json:"ids"`
	}

	var response struct {
		Data map[int]int `json:"data"`
	}

	res, err := ts.client.R().
		SetHeader("content-type", "application/json").
		SetBody(requestBody{Ids: todoIds}).
		SetResult(&response).
		Post(fmt.Sprintf("%s/%s", ts.serviceUrl, "rpc/get_todo_total_react")) // todo: use var to handle this api name

	log.Println(fmt.Sprintf("%s/%s", ts.serviceUrl, "rpc/get_todo_total_react"))

	if err != nil {
		return nil, err
	}

	if !res.IsSuccess() {
		//log.Println(res.RawResponse)
		return nil, errors.New(res.String())
	}

	return response.Data, nil
}
