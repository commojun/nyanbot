package main

import (
	"fmt"
	"log"

	"github.com/commojun/nyanbot/masterdata/key_value"
)

func main() {
	kvs, err := key_value.New()
	if err != nil {
		log.Fatal(err)
	}

	kvs.LoadKVsFromSheet()

	fmt.Println(kvs)
	fmt.Println(kvs.Rooms)

}
