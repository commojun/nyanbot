package main

import (
	"fmt"
	"log"

	"github.com/commojun/nyanbot/app/redis"
	"github.com/commojun/nyanbot/masterdata/key_value"
)

func main() {
	kvs, err := key_value.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*kvs)

	redisClient := redis.NewClient()
	roomIds, err := redisClient.HGetAll(key_value.Room).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(roomIds)
}
