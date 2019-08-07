package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/commojun/nyanbot/masterdata/table"
)

func main() {
	ts, err := table.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("----------構造体からJSONへ----------")
	jsonBytes, err := json.Marshal(ts.Alarms)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonBytes))

	fmt.Printf("\n\n")

	fmt.Println("----------Redisから復元する----------")
	fmt.Println("alarm")
	alms, err := table.Alarms()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*alms)

	fmt.Println("anniversary")
	anvs, err := table.Anniversaries()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*anvs)
}
