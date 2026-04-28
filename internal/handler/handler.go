package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/commojun/nyanbot/internal/config"
)

type API struct {
	Name string
	Get  func(context.Context, *http.Request, *Response) error
	Post func(context.Context, *http.Request, *Response) error
}

type Response struct {
	Status  int
	Message string
	Writer  *http.ResponseWriter
}

func newAPIs(cfg config.Config) ([]API, error) {
	return []API{
		test,
		makeMessageAPI(cfg),
		makeLineHookAPI(cfg),
	}, nil
}

func (api *API) MakeHundleFunc() (string, func(http.ResponseWriter, *http.Request)) {
	return api.Name, func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		res := Response{
			Status:  http.StatusOK,
			Message: "OK",
			Writer:  &w,
		}
		var err error

		if req.Method == http.MethodPost && api.Post != nil {
			err = api.Post(ctx, req, &res)
		} else if req.Method == http.MethodGet && api.Get != nil {
			err = api.Get(ctx, req, &res)
		} else {
			// リクエストされたメソッドが用意されていない
			res.Status = http.StatusMethodNotAllowed
			err = fmt.Errorf("Method not allowed")
		}

		if err != nil {
			// エラーを作る
			res.Message = fmt.Sprintf("[API:%s] %d error: %s", api.Name, res.Status, err)
			log.Println(res.Message)
		} else {
			log.Printf("[API:%s] %d OK", api.Name, res.Status)
		}

		w.WriteHeader(res.Status)
		w.Write([]byte(res.Message))

		return
	}
}

func parseJSONRequest(req *http.Request, i interface{}) error {
	// POSTでないと受け付けない
	if req.Method != "POST" {
		return fmt.Errorf("Method is not POST")
	}

	// JSONでないと受け付けない
	if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		return fmt.Errorf("Content-Type is not JSON")
	}

	// Content-Lengthを取得
	length, err := strconv.Atoi(req.Header.Get("Content-Length"))
	if err != nil {
		return fmt.Errorf("Could not get Content-Length")
	}

	// リクエストBodyを取得
	body := make([]byte, length)
	length, err = req.Body.Read(body)
	if err != nil && err != io.EOF {
		return fmt.Errorf("Read body failed")
	}

	// jsonを構造体に当てはめる
	err = json.Unmarshal(body[:length], i)
	if err != nil {
		return fmt.Errorf("JSON unmarshal failed")
	}

	return nil
}
