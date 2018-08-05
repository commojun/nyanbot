package main

import (
	"log"

	"github.com/junpooooow/nyanbot"
)

func main() {
	err := nyanbot.SendPushMessage()
	if err != nil {
		log.Fatal(err)
	}
}
