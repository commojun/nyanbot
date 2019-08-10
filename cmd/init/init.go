package main

import (
	"log"
	"time"

	"github.com/Songmu/retry"
	"github.com/commojun/nyanbot/app/redis"
	"github.com/commojun/nyanbot/masterdata/key_value"
	"github.com/commojun/nyanbot/masterdata/table"
)

func main() {
	// Redisに接続できるか
	rc := redis.NewClient()
	err := retry.Retry(10, 10*time.Second, func() error {
		log.Println("attempt to connect to redis")
		err := rc.Keys("*").Err()
		return err
	})
	if err != nil {
		log.Println("failed to connect to redis")
		log.Fatal(err)
	}
	log.Println("redis connection ok")

	log.Println("Table initialize")
	_, err = table.Initialize()
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
