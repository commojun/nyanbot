package main

import (
	"log"

	"github.com/commojun/nyanbot/masterdata/key_value"
	"github.com/commojun/nyanbot/masterdata/table"
)

func main() {
	log.Println("Table initialize")
	_, err := table.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("KeyValue initialize")
	_, err = key_value.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Initialize done.")
}
