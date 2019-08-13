package server

import (
	"log"

	"net/http"

	"github.com/commojun/nyanbot/api"
)

type Server struct {
}

func New() (*Server, error) {
	return &Server{}, nil
}

func (server *Server) Start() error {
	apis := api.New()
	for _, api := range apis {
		http.HandleFunc(api.MakeHundleFunc())
		log.Printf("Registered API: %s", api.Name)
	}

	port := "8999"
	log.Printf("Start server port:%s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return err
	}
	return nil
}
