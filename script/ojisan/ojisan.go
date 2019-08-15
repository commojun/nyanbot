package main

import (
	"fmt"
	"log"

	"github.com/commojun/nyanbot/app/ojisan"
)

func main() {
	name := "じゅん"
	emojiNum := 4
	level := 2

	ojisan := ojisan.New(name, emojiNum, level)

	msg, err := ojisan.Generate()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(msg)
}
