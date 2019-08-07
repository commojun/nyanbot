package main

import (
	"log"

	"github.com/commojun/nyanbot/app/sheet"
)

func main() {
	sheet, err := sheet.New()
	if err != nil {
		log.Fatal(err)
	}

	err = sheet.Load()
	if err != nil {
		log.Fatal(err)
	}
}
