package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type API struct {
	Name string
	Get  func(*http.Request, *Response) error
	Post func(*http.Request, *Response) error
}

type Response struct {
	Status  int
	Message string
	Writer  *http.ResponseWriter
}

func New() ([]API, error) {
	return []API{
		test,
		message,
	}, nil
}

func (api *API) MakeHundleFunc() (string, func(http.ResponseWriter, *http.Request)) {
	return api.Name, func(w http.ResponseWriter, req *http.Request) {
		res := Response{
			Status:  http.StatusOK,
			Message: "OK",
			Writer:  &w,
		}
		var err error

		if req.Method == http.MethodPost && api.Post != nil {
			err = api.Post(req, &res)
		} else if req.Method == http.MethodGet && api.Get != nil {
			err = api.Get(req, &res)
		} else {
			// リクエストされたメソッドが用意されていない
			res.Status = http.StatusMethodNotAllowed
			err = fmt.Errorf("Method not allowed")
		}

		if err != nil {
			// エラーを作る
			res.Message = fmt.Sprintf("[API:%s] %d error: %s, request: %s", api.Name, res.Status, err, *req)
			log.Println(res.Message)
		} else {
			log.Printf("[API:%s] %d OK, request: %s", api.Name, res.Status, *req)
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
	if req.Header.Get("Content-Type") != "application/json" {
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
