package main

import (
	"log"

	"github.com/commojun/nyanbot/app/echo"
)

func main() {
	echo, err := echo.New()
	if err != nil {
		log.Fatal(err)
	}

	err = echo.StartServer()
	if err != nil {
		log.Fatal(err)
	}
}
