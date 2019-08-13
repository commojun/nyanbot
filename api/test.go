package api

import "net/http"

var (
	test = API{
		Name: "/test",
		Get:  getTest,
	}
)

func getTest(req *http.Request, res *Response) error {
	res.Message = "this is test"
	return nil
}
