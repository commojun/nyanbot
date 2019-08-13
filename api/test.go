package api

import "net/http"

var (
	Test = API{
		Name: "/test",
		Get:  GetTest,
		Post: nil,
	}
)

func GetTest(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("this is test"))
	return
}
