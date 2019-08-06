package redis

import (
	"github.com/commojun/nyanbot/app/constant"
	origin "github.com/go-redis/redis"
)

func NewClient() *origin.Client {
	client := origin.NewClient(&origin.Options{
		Addr:     constant.RedisHost,
		Password: "",
		DB:       0,
	})

	return client
}
