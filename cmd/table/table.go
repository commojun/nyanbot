package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/commojun/nyanbot/masterdata/table"
)

func main() {
	ts, err := table.New()
	if err != nil {
		log.Fatal(err)
	}

	ts.LoadTablesFromSheet()

	jsonBytes, err := json.Marshal(ts.Alarms)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("----------構造体からJSONへ----------")
	fmt.Println(string(jsonBytes))

	fmt.Printf("\n\n")

	fmt.Println("----------JSONから構造体へ----------")
	alms := []table.Alarm{}
	json.Unmarshal(jsonBytes, &alms)
	fmt.Println(alms)

}
