package main

import (
	"log"

	"github.com/commojun/nyanbot"
)

func main() {
	server, err := nyanbot.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
