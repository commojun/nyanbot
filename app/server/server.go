package server

import (
	"log"
	"strconv"

	"net/http"

	"github.com/commojun/nyanbot/api"
	"github.com/commojun/nyanbot/constant"
)

type Server struct {
	APIs []api.API
	Port int
}

func New() (*Server, error) {
	port, err := strconv.Atoi(constant.ServerPort)
	if err != nil {
		if constant.ServerPort == "" {
			port = constant.DefaultServerPort
		} else {
			return nil, err
		}
	}

	apis, err := api.New()
	return &Server{
		APIs: apis,
		Port: port,
	}, nil
}

func (server *Server) Start() error {
	for _, api := range server.APIs {
		http.HandleFunc(api.MakeHundleFunc())
		log.Printf("Registered API: %s", api.Name)
	}

	port := strconv.Itoa(server.Port)
	log.Printf("Start server port:%s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return err
	}
	return nil
}
