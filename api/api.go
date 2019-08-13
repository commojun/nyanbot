package api

import "net/http"

type API struct {
	Name string
	Get  func(http.ResponseWriter, *http.Request)
	Post func(http.ResponseWriter, *http.Request)
}

var (
	APIs = []*API{
		&API{
			Name: "/hoge",
			Get:  methodNotAllowed,
			Post: methodNotAllowed,
		},
		&Test,
	}
)

func (api *API) MakeHundleFunc() (string, func(http.ResponseWriter, *http.Request)) {
	return api.Name, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost && api.Post != nil {
			api.Post(w, req)
		} else if req.Method == http.MethodGet && api.Get != nil {
			api.Get(w, req)
		} else {
			methodNotAllowed(w, req)
		}
		return
	}
}

func methodNotAllowed(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed) // 405
	w.Write([]byte("Method not allowed"))
	return
}
