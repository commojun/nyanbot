package nyanbot

import "github.com/commojun/nyanbot/app/server"

func NewServer() (*server.Server, error) {
	s, err := server.New()
	if err != nil {
		return &server.Server{}, err
	}
	return s, nil
}
