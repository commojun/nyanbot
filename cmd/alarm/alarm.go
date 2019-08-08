package main

import (
	"log"

	"github.com/commojun/nyanbot/app/alarm"
)

func main() {
	alm, err := alarm.New()
	if err != nil {
		log.Fatal(err)
	}

	err = alm.Run()
	if err != nil {
		log.Fatal(err)
	}
}
