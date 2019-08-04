package main

import (
	"log"

	"github.com/commojun/nyanbot/table"
)

func main() {
	t, err := table.New()
	if err != nil {
		log.Fatal(err)
	}

	t.LoadFromSheet()
}
