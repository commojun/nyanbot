package main

import (
	"fmt"
	"log"

	"github.com/commojun/nyanbot/app/redis"
)

func main() {
	client := redis.NewClient()

	val, err := client.Set("key", "value", 0).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(val)

	val, err = client.Get("key").Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("key", val)
}
