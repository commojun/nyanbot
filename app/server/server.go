package server

import (
	"log"
	"strconv"

	"net/http"

	"github.com/commojun/nyanbot/api"
	"github.com/commojun/nyanbot/config"
)

type Server struct {
	APIs []api.API
	Port int
}

func New(cfg config.Config) (*Server, error) {
	apis, err := api.New(cfg)
	return &Server{
		APIs: apis,
		Port: cfg.ServerPort,
	}, err
}

func (server *Server) Start() error {
	for _, api := range server.APIs {
		newApi := api
		http.HandleFunc(newApi.MakeHundleFunc())
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
