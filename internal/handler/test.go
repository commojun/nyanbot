package handler

import (
	"context"
	"net/http"
)

var (
	test = API{
		Name: "/test",
		Get:  getTest,
	}
)

func getTest(ctx context.Context, req *http.Request, res *Response) error {
	res.Message = "this is test"
	return nil
}
