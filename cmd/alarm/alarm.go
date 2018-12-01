package main

import (
	"log"

	"github.com/commojun/nyanbot/app/alarm"
)

func main() {
	alm, err := alarm.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = alm.Send()
	if err != nil {
		log.Fatal(err)
	}
}
