package main

import (
	"log"

	"github.com/commojun/nyanbot/app/anniversary"
)

func main() {
	am, err := anniversary.New()
	if err != nil {
		log.Fatal(err)
	}

	err = am.Run()
	if err != nil {
		log.Fatal(err)
	}
}
