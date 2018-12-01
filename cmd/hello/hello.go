package main

import (
	"log"

	"github.com/commojun/nyanbot/app/hello"
)

func main() {
	hello, err := hello.New()
	if err != nil {
		log.Fatal(err)
	}

	err = hello.Say()
	if err != nil {
		log.Fatal(err)
	}
}
